package handlers

import (
	"net/http"

	"conduit/internal/repository"
	"conduit/internal/services"

	"github.com/gin-gonic/gin"
)

func GetPerformanceMetrics(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		metrics, err := services.CalculatePerformanceMetrics(repo, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate performance metrics"})
			return
		}

		c.JSON(http.StatusOK, metrics)
	}
}

func GetRiskMetrics(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		metrics, err := services.CalculateRiskMetrics(repo, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate risk metrics"})
			return
		}

		c.JSON(http.StatusOK, metrics)
	}
}
