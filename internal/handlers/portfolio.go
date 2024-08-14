package handlers

import (
	"net/http"

	"conduit/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func GetPortfolio(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		portfolio, err := repo.GetPortfolio(userID.(uint))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
			return
		}

		c.JSON(http.StatusOK, portfolio)
	}
}
