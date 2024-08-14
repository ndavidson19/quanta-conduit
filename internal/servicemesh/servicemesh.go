package servicemesh

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ServiceMesh struct {
	services map[string][]string
	mu       sync.RWMutex
	logger   *zap.Logger
}

func NewServiceMesh(logger *zap.Logger) *ServiceMesh {
	return &ServiceMesh{
		services: make(map[string][]string),
		logger:   logger,
	}
}

func (sm *ServiceMesh) RegisterService(name, url string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.services[name]; !exists {
		sm.services[name] = []string{}
	}
	sm.services[name] = append(sm.services[name], url)
	sm.logger.Info("Service registered", zap.String("name", name), zap.String("url", url))
}

func (sm *ServiceMesh) DeregisterService(name, url string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if urls, exists := sm.services[name]; exists {
		for i, u := range urls {
			if u == url {
				sm.services[name] = append(urls[:i], urls[i+1:]...)
				sm.logger.Info("Service deregistered", zap.String("name", name), zap.String("url", url))
				return
			}
		}
	}
}

func (sm *ServiceMesh) GetService(name string) ([]string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	urls, exists := sm.services[name]
	return urls, exists
}

func (sm *ServiceMesh) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	json.NewEncoder(w).Encode(sm.services)
}

func (sm *ServiceMesh) HealthCheck() {
	for {
		sm.mu.RLock()
		for name, urls := range sm.services {
			for _, url := range urls {
				go func(name, url string) {
					if !isServiceAlive(url) {
						sm.DeregisterService(name, url)
					}
				}(name, url)
			}
		}
		sm.mu.RUnlock()
		time.Sleep(30 * time.Second)
	}
}

func isServiceAlive(url string) bool {
	resp, err := http.Get(url + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}