package main

import (
    "log"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "conduit/internal/api"
    "conduit/internal/handlers"
    "conduit/internal/store"
	"conduit/pkg/utils"
)

func main() {
    // Initialize shared utilities
    utils.InitLogger()
    err := utils.LoadConfig("./config/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

	dbURL := "postgres://postgres:my_postgres_password@localhost:5432/my_postgres_db?sslmode=disable"  
	if val, ok := utils.GetConfig("database_url").(string); ok && val != "" {
        dbURL = val
    }
    log.Printf("Initializing store with DATABASE_URL: %s", dbURL)

    storeConfig := store.Config{
        PostgresURL: dbURL,
    }
    store, err := store.NewStore(storeConfig)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer store.Close()

    // Create a new Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Create an instance of our handler
    h := handlers.NewHandler(store)

    // Register our handler to implement the generated interface
    api.RegisterHandlers(e, h)

    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}