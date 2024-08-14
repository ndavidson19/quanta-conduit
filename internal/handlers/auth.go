package handlers

import (
	"net/http"

	"conduit/internal/auth"
	"conduit/internal/repository/postgres"
	"conduit/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Register(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashedPassword

		if err := repo.CreateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func Login(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetUserByUsername(loginData.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if !utils.CheckPasswordHash(loginData.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		accessToken, refreshToken, err := auth.GenerateTokenPair(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshToken(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := c.GetHeader("Refresh-Token")
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
			return
		}

		userID, err := auth.ValidateRefreshToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		user, err := repo.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}

		newAccessToken, newRefreshToken, err := auth.GenerateTokenPair(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	}
}
