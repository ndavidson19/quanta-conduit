package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type LoadBalancingAlgorithm int

const (
	RoundRobin LoadBalancingAlgorithm = iota
	LeastConnections
	IPHash
)

type AdvancedLoadBalancer struct {
	backends       []*Backend
	algorithm      LoadBalancingAlgorithm
	connectionCount map[*Backend]*int64
	mutex          sync.RWMutex
	logger         *zap.Logger
	metrics        *Metrics
}

type Metrics struct {
	RequestCount    *prometheus.CounterVec
	ResponseTime    *prometheus.HistogramVec
	ActiveConnections *prometheus.GaugeVec
}

func NewAdvancedLoadBalancer(backends []string, algorithm LoadBalancingAlgorithm, logger *zap.Logger) (*AdvancedLoadBalancer, error) {
	lb := &AdvancedLoadBalancer{
		algorithm:      algorithm,
		connectionCount: make(map[*Backend]*int64),
		logger:         logger,
		metrics:        initMetrics(),
	}

	for _, backend := range backends {
		b, err := NewBackend(backend)
		if err != nil {
			return nil, err
		}
		lb.backends = append(lb.backends, b)
		count := int64(0)
		lb.connectionCount[b] = &count
	}

	return lb, nil
}

func (lb *AdvancedLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	peer := lb.getNextPeer(r)
	if peer != nil {
		start := time.Now()
		atomic.AddInt64(lb.connectionCount[peer], 1)
		lb.metrics.ActiveConnections.WithLabelValues(peer.URL.Host).Inc()

		peer.ReverseProxy.ServeHTTP(w, r)

		atomic.AddInt64(lb.connectionCount[peer], -1)
		lb.metrics.ActiveConnections.WithLabelValues(peer.URL.Host).Dec()
		lb.metrics.RequestCount.WithLabelValues(peer.URL.Host).Inc()
		lb.metrics.ResponseTime.WithLabelValues(peer.URL.Host).Observe(time.Since(start).Seconds())
	} else {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}

func (lb *AdvancedLoadBalancer) getNextPeer(r *http.Request) *Backend {
	switch lb.algorithm {
	case RoundRobin:
		return lb.roundRobin()
	case LeastConnections:
		return lb.leastConnections()
	case IPHash:
		return lb.ipHash(r)
	default:
		return lb.roundRobin()
	}
}

func (lb *AdvancedLoadBalancer) roundRobin() *Backend {
	// Implementation remains the same as before
}

func (lb *AdvancedLoadBalancer) leastConnections() *Backend {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	var leastConn *Backend
	minConn := int64(^uint64(0) >> 1) // Max int64 value

	for backend, count := range lb.connectionCount {
		if backend.IsAlive() && atomic.LoadInt64(count) < minConn {
			leastConn = backend
			minConn = atomic.LoadInt64(count)
		}
	}

	return leastConn
}

func (lb *AdvancedLoadBalancer) ipHash(r *http.Request) *Backend {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	ip := getIP(r)
	hash := hashIP(ip)
	index := hash % uint32(len(lb.backends))

	return lb.backends[index]
}

// Helper functions (getIP, hashIP) to be implemented