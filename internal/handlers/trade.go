package handlers

import (
	"net/http"

	"conduit/internal/models"
	"conduit/internal/repository"

	"github.com/gin-gonic/gin"
)

func PlaceTrade(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var trade models.Trade
		if err := c.ShouldBindJSON(&trade); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("userID")
		trade.UserID = userID.(uint)

		if err := repo.CreateTrade(&trade); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place trade"})
			return
		}

		c.JSON(http.StatusCreated, trade)
	}
}

func ListTrades(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		trades, err := repo.ListTradesByUser(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades"})
			return
		}

		c.JSON(http.StatusOK, trades)
	}
}

// Implement GetTrade and CancelTrade handlers
