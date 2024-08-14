package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"fmt"
	"testing"
	"time"

	"conduit/internal/config"
	"conduit/internal/server"
	"conduit/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRepository is a mock implementation of the repository.Repository interface
type MockRepository struct {
	mock.Mock
}

// Implement all methods of the repository.Repository interface here...

func getServerAddress() string {
	host := os.Getenv("TEST_SERVER_HOST")
	port := os.Getenv("TEST_SERVER_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}

func TestServerInitialization(t *testing.T) {
	cfg := &config.Config{
		ServerPort: "8080",
		// Set other necessary config fields
	}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	assert.NotNil(t, s, "Server should not be nil")
	assert.NotNil(t, s.Router(), "Router should not be nil")
}

func TestMiddlewareSetup(t *testing.T) {
	cfg := &config.Config{ServerPort: "8080"}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	// Test CORS middleware
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/", nil)
	s.Router().ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))

	// Test other middleware as needed...
}

func TestRouteSetup(t *testing.T) {
	cfg := &config.Config{ServerPort: "8080"}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	routes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/auth/register"},
		{"POST", "/api/v1/auth/login"},
		{"GET", "/api/v1/users/me"},
		// Add all other routes here...
	}

	for _, route := range routes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(route.method, route.path, nil)
		s.Router().ServeHTTP(w, req)

		assert.NotEqual(t, 404, w.Code, "Route %s %s should exist", route.method, route.path)
	}
}

func TestServerShutdown(t *testing.T) {
	cfg := &config.Config{ServerPort: "8080"}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	go func() {
		err := s.Run()
		assert.NoError(t, err, "Server should run without error")
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	assert.NoError(t, err, "Server should shut down gracefully")
}

func TestConcurrentRequests(t *testing.T) {
	cfg := &config.Config{ServerPort: "8080"}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	go s.Run()
	defer s.Shutdown(context.Background())

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	concurrentRequests := 100
	done := make(chan bool)

	serverAddress := getServerAddress()

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			resp, err := http.Get(serverAddress + "/api/v1/users/me")
			assert.NoError(t, err, "Request should not error")
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Unauthenticated request should return 401")
			done <- true
		}()
	}

	for i := 0; i < concurrentRequests; i++ {
		<-done
	}
}

func TestSecurityHeaders(t *testing.T) {
	cfg := &config.Config{ServerPort: "8081"}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	go s.Run()
	defer s.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	serverAddress := getServerAddress()

	resp, err := http.Get(serverAddress + "/api/v1/users/me")
	assert.NoError(t, err, "Request should not error")

	expectedHeaders := map[string]string{
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"X-XSS-Protection":          "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for header, expectedValue := range expectedHeaders {
		assert.Equal(t, expectedValue, resp.Header.Get(header), fmt.Sprintf("%s header should be set correctly", header))
	}
}

func TestRateLimiting(t *testing.T) {
	cfg := &config.Config{ServerPort: "8082"}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	go s.Run()
	defer s.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	client := &http.Client{}
	serverAddress := getServerAddress()
	url := serverAddress + "/api/v1/users/me"

	for i := 0; i < 100; i++ {
		resp, err := client.Get(url)
		assert.NoError(t, err, "Request should not error")
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			// Rate limit hit, test passed
			return
		}
	}

	t.Error("Rate limiting not triggered after 100 requests")
}

func TestLoadTesting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	cfg := &config.Config{ServerPort: "8083"}
	logger, _ := zap.NewDevelopment()
	repo := new(MockRepository)

	s := server.NewServer(cfg, logger, repo)

	go s.Run()
	defer s.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	serverAddress := getServerAddress()
	url := serverAddress + "/api/v1/users/me"
	concurrency := 50
	totalRequests := 1000

	var wg sync.WaitGroup
	results := make(chan int, totalRequests)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrency; j++ {
				resp, err := http.Get(url)
				if err != nil {
					t.Errorf("Request error: %v", err)
					results <- -1
				} else {
					results <- resp.StatusCode
					resp.Body.Close()
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	statusCodes := make(map[int]int)
	for code := range results {
		statusCodes[code]++
	}

	t.Logf("Load test results: %v", statusCodes)
	assert.Equal(t, totalRequests, statusCodes[http.StatusUnauthorized], "All requests should return 401 Unauthorized")
}
