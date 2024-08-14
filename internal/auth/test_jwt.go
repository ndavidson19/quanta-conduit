package auth

import (
	"conduit/internal/models"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	// Set up
	os.Setenv("JWT_SECRET", "test_secret")
	user := &models.User{ID: 1}

	// Test
	token, err := GenerateToken(user)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "quant-trading-microservice", claims.Issuer)
}

func TestValidateToken(t *testing.T) {
	// Set up
	os.Setenv("JWT_SECRET", "test_secret")
	user := &models.User{ID: 1}
	token, _ := GenerateToken(user)

	// Test
	claims, err := ValidateToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, uint(1), claims.UserID)
}

func TestGenerateTokenPair(t *testing.T) {
	// Set up
	os.Setenv("JWT_SECRET", "test_secret")
	user := &models.User{ID: 1}

	// Test
	accessToken, refreshToken, err := GenerateTokenPair(user)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Validate access token
	accessClaims, err := ValidateToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), accessClaims.UserID)

	// Validate refresh token
	userID, err := ValidateRefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), userID)
}

func TestValidateRefreshToken(t *testing.T) {
	// Set up
	os.Setenv("JWT_SECRET", "test_secret")
	user := &models.User{ID: 1}
	_, refreshToken, _ := GenerateTokenPair(user)

	// Test
	userID, err := ValidateRefreshToken(refreshToken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(1), userID)
}