package loadbalancer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewAdvancedLoadBalancer(t *testing.T) {
	// Test
	backends := []string{"http://localhost:8081", "http://localhost:8082"}
	lb, err := NewAdvancedLoadBalancer(backends, RoundRobin, zap.NewNop())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, lb)
	assert.Len(t, lb.backends, 2)
	assert.Equal(t, RoundRobin, lb.algorithm)
}

func TestRoundRobin(t *testing.T) {
	// Set up
	backends := []string{"http://localhost:8081", "http://localhost:8082"}
	lb, _ := NewAdvancedLoadBalancer(backends, RoundRobin, zap.NewNop())

	// Test
	backend1 := lb.roundRobin()
	backend2 := lb.roundRobin()
	backend3 := lb.roundRobin()

	// Assert
	assert.NotEqual(t, backend1, backend2)
	assert.Equal(t, backend1, backend3)
}

func TestLeastConnections(t *testing.T) {
	// Set up
	backends := []string{"http://localhost:8081", "http://localhost:8082"}
	lb, _ := NewAdvancedLoadBalancer(backends, LeastConnections, zap.NewNop())

	// Simulate connections
	*lb.connectionCount[lb.backends[0]] = 5
	*lb.connectionCount[lb.backends[1]] = 2

	// Test
	backend := lb.leastConnections()

	// Assert
	assert.Equal(t, lb.backends[1], backend)
}

func TestIPHash(t *testing.T) {
	// Set up
	backends := []string{"http://localhost:8081", "http://localhost:8082"}
	lb, _ := NewAdvancedLoadBalancer(backends, IPHash, zap.NewNop())

	// Test
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "192.168.1.2:12345"

	backend1 := lb.ipHash(req1)
	backend2 := lb.ipHash(req2)
	backend3 := lb.ipHash(req1)

	// Assert
	assert.NotEqual(t, backend1, backend2)
	assert.Equal(t, backend1, backend3)
}

func TestServeHTTP(t *testing.T) {
	// Set up
	backends := []string{"http://localhost:8081", "http://localhost:8082"}
	lb, _ := NewAdvancedLoadBalancer(backends, RoundRobin, zap.NewNop())

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Replace the first backend with the test server
	lb.backends[0].URL = ts.URL

	// Test
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	lb.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}