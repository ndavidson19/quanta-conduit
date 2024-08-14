package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	backends []*Backend
	current  uint64
	logger   *zap.Logger
}

func NewLoadBalancer(backends []string, logger *zap.Logger) (*LoadBalancer, error) {
	var be []*Backend

	for _, b := range backends {
		url, err := url.Parse(b)
		if err != nil {
			return nil, err
		}

		backend := &Backend{
			URL:          url,
			Alive:        true,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		}
		be = append(be, backend)
	}

	return &LoadBalancer{
		backends: be,
		logger:   logger,
	}, nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	peer := lb.getNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}

func (lb *LoadBalancer) getNextPeer() *Backend {
	next := int(atomic.AddUint64(&lb.current, uint64(1)) % uint64(len(lb.backends)))
	l := len(lb.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(lb.backends)
		if lb.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&lb.current, uint64(idx))
			}
			return lb.backends[idx]
		}
	}
	return nil
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

func (lb *LoadBalancer) HealthCheck() {
	for _, b := range lb.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		lb.logger.Info("Backend status", zap.String("url", b.URL.String()), zap.String("status", status))
	}
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}