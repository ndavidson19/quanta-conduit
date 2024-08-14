package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"conduit/internal/handlers"
	"conduit/internal/repository/postgres"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateStrategy(t *testing.T) {
	// Setup test database and repository
	repo := setupTestRepo()

	// Setup router and handler
	router := gin.New()
	router.POST("/strategy", handlers.CreateStrategy(repo))

	// Create test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/strategy", strings.NewReader(`{"name":"Test Strategy","parameters":"{}"}`))
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, 201, w.Code)
}

func setupTestRepo() *postgres.Repository {
	// Implement test repository setup
	// This could involve setting up an in-memory database or mocking the repository
	return &postgres.Repository{}
}
