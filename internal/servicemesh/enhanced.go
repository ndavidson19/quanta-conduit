package servicemesh

import (
	"context"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type EnhancedServiceMesh struct {
	etcdClient *clientv3.Client
	services   map[string][]string
	mutex      sync.RWMutex
	logger     *zap.Logger
}

func NewEnhancedServiceMesh(endpoints []string, logger *zap.Logger) (*EnhancedServiceMesh, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	sm := &EnhancedServiceMesh{
		etcdClient: cli,
		services:   make(map[string][]string),
		logger:     logger,
	}

	go sm.watchServices()

	return sm, nil
}

func (sm *EnhancedServiceMesh) RegisterService(ctx context.Context, name, url string, ttl int64) error {
	key := "/services/" + name + "/" + url
	lease, err := sm.etcdClient.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	_, err = sm.etcdClient.Put(ctx, key, url, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	// Keep the lease alive
	keepAliveCh, err := sm.etcdClient.KeepAlive(ctx, lease.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-keepAliveCh:
				// Lease kept alive successfully
			case <-ctx.Done():
				sm.etcdClient.Revoke(context.Background(), lease.ID)
				return
			}
		}
	}()

	sm.logger.Info("Service registered", zap.String("name", name), zap.String("url", url))
	return nil
}

func (sm *EnhancedServiceMesh) DeregisterService(ctx context.Context, name, url string) error {
	key := "/services/" + name + "/" + url
	_, err := sm.etcdClient.Delete(ctx, key)
	if err != nil {
		return err
	}

	sm.logger.Info("Service deregistered", zap.String("name", name), zap.String("url", url))
	return nil
}

func (sm *EnhancedServiceMesh) GetService(ctx context.Context, name string) ([]string, error) {
	resp, err := sm.etcdClient.Get(ctx, "/services/"+name+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, kv := range resp.Kvs {
		urls = append(urls, string(kv.Value))
	}

	return urls, nil
}

func (sm *EnhancedServiceMesh) watchServices() {
	watchChan := sm.etcdClient.Watch(context.Background(), "/services/", clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			sm.handleServiceEvent(event)
		}
	}
}

func (sm *EnhancedServiceMesh) handleServiceEvent(event *clientv3.Event) {
	// Implementation to handle service updates
}