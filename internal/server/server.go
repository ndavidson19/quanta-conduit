package server

import (
	"context"
	"net/http"

	"conduit/internal/config"
	"conduit/internal/handlers"
	"conduit/internal/kafka"
	"conduit/internal/middleware"
	"conduit/internal/repository"
	"conduit/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	router    *gin.Engine
	config    *config.Config
	logger    *zap.Logger
	repo      repository.Repository
	srv       *http.Server
	wsManager *websocket.Manager
	kafka     *kafka.Consumer
}

func NewServer(cfg *config.Config, logger *zap.Logger, repo repository.Repository) *Server {
	router := gin.New()

	wsManager := websocket.NewManager(logger)

	kafkaConsumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaGroupID,
		cfg.KafkaTopics,
		logger,
		wsManager,
	)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", zap.Error(err))
	}

	s := &Server{
		router:    router,
		config:    cfg,
		logger:    logger,
		repo:      repo,
		wsManager: wsManager,
		kafka:     kafkaConsumer,
	}

	s.setupMiddleware()
	s.setupRoutes()

	s.srv = &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: s.router,
	}

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Logger(s.logger))
	s.router.Use(middleware.Recovery(s.logger))
	s.router.Use(middleware.CORS())
	s.router.Use(middleware.SecurityHeaders())
}

func (s *Server) setupRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register(s.repo))
			auth.POST("/login", handlers.Login(s.repo))
			auth.POST("/refresh", handlers.RefreshToken(s.repo))
			auth.POST("/logout", middleware.Auth(), handlers.Logout(s.repo))
		}

		user := v1.Group("/users").Use(middleware.Auth())
		{
			user.GET("/me", handlers.GetCurrentUser(s.repo))
			user.PUT("/me", handlers.UpdateUser(s.repo))
			user.GET("/portfolio", handlers.GetUserPortfolio(s.repo))
		}

		strategy := v1.Group("/strategies").Use(middleware.Auth())
		{
			strategy.POST("/", handlers.CreateStrategy(s.repo))
			strategy.GET("/", handlers.ListStrategies(s.repo))
			strategy.GET("/:id", handlers.GetStrategy(s.repo))
			strategy.PUT("/:id", handlers.UpdateStrategy(s.repo))
			strategy.DELETE("/:id", handlers.DeleteStrategy(s.repo))
			strategy.POST("/:id/backtest", handlers.BacktestStrategy(s.repo))
			strategy.POST("/:id/deploy", handlers.DeployStrategy(s.repo))
		}

		datasource := v1.Group("/datasources").Use(middleware.Auth())
		{
			datasource.GET("/", handlers.ListDataSources(s.repo))
			datasource.POST("/subscribe/:id", handlers.SubscribeToDataSource(s.repo))
			datasource.POST("/unsubscribe/:id", handlers.UnsubscribeFromDataSource(s.repo))
		}

		trade := v1.Group("/trades").Use(middleware.Auth())
		{
			trade.POST("/", handlers.PlaceTrade(s.repo))
			trade.GET("/", handlers.ListTrades(s.repo))
			trade.GET("/:id", handlers.GetTrade(s.repo))
			trade.DELETE("/:id", handlers.CancelTrade(s.repo))
		}

		analytics := v1.Group("/analytics").Use(middleware.Auth())
		{
			analytics.GET("/performance", handlers.GetPerformanceMetrics(s.repo))
			analytics.GET("/risk", handlers.GetRiskMetrics(s.repo))
		}

		s.router.GET("/ws", s.handleWebSocket)
	}
}

func (s *Server) handleWebSocket(c *gin.Context) {
	s.wsManager.AuthenticateAndServeWS(c.Writer, c.Request, s.logger)
}

func (s *Server) Run() error {
	go s.wsManager.Run()
	go s.kafka.ConsumeWithRetry(context.Background())

	s.logger.Info("Starting server", zap.String("port", s.config.ServerPort))
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	return s.srv.Shutdown(ctx)
}
