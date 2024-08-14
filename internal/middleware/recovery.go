package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                logger.Error("panic recovered",
                    zap.Any("error", err),
                    zap.String("request", c.Request.URL.Path),
                )
                c.AbortWithStatus(http.StatusInternalServerError)
            }
        }()
        c.Next()
    }
}