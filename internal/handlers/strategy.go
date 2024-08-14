package handlers

import (
	"net/http"
	"strconv"

	"conduit/internal/models"
	"conduit/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func CreateStrategy(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var strategy models.Strategy
		if err := c.ShouldBindJSON(&strategy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get user ID from the authenticated context
		userID, _ := c.Get("userID")
		strategy.UserID = userID.(uint)

		if err := repo.CreateStrategy(&strategy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create strategy"})
			return
		}

		c.JSON(http.StatusCreated, strategy)
	}
}

func GetStrategy(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		strategy, err := repo.GetStrategy(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Strategy not found"})
			return
		}

		c.JSON(http.StatusOK, strategy)
	}
}

func UpdateStrategy(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var strategy models.Strategy
		if err := c.ShouldBindJSON(&strategy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		strategy.ID = uint(id)
		if err := repo.UpdateStrategy(&strategy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update strategy"})
			return
		}

		c.JSON(http.StatusOK, strategy)
	}
}
