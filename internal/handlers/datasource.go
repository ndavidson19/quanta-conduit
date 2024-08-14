package handlers

import (
	"net/http"
	"strconv"

	"conduit/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func ListDataSources(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		dataSources, err := repo.ListDataSources()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data sources"})
			return
		}

		c.JSON(http.StatusOK, dataSources)
	}
}

func SubscribeToDataSource(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		dataSourceID, _ := strconv.Atoi(c.Param("id"))

		if err := repo.SubscribeToDataSource(userID.(uint), uint(dataSourceID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe to data source"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully"})
	}
}
