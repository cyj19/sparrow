/**
 * @Author: cyj19
 * @Date: 2022/3/10 18:51
 */

package discovery

import (
	"errors"
	"github.com/cyj19/sparrow/balance"
	"github.com/cyj19/sparrow/registry"
	"sync"
)

type Discovery interface {
	Refresh() error                              // 从注册中心更新服务列表
	Update(servers []*registry.ServerItem) error // 手动更新服务列表
	Get() (*registry.ServerItem, error)
	GetAll() ([]*registry.ServerItem, error)
}

type SimpleDiscovery struct {
	mu      *sync.Mutex
	servers []*registry.ServerItem
	lb      balance.LoadBalancing
}

var _ Discovery = (*SimpleDiscovery)(nil)

func NewSimpleDiscovery(lb balance.LoadBalancing) *SimpleDiscovery {
	return &SimpleDiscovery{
		mu:      new(sync.Mutex),
		servers: make([]*registry.ServerItem, 0),
		lb:      lb,
	}
}

func (d *SimpleDiscovery) Refresh() error {
	return nil
}

func (d *SimpleDiscovery) Update(servers []*registry.ServerItem) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

func (d *SimpleDiscovery) Get() (*registry.ServerItem, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := len(d.servers)
	if n == 0 {
		return nil, errors.New("rpc discovery: no available servers")
	}
	idx := d.lb.GetModeResult(n)
	s := d.servers[idx]
	return s, nil
}

func (d *SimpleDiscovery) GetAll() ([]*registry.ServerItem, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	servers := make([]*registry.ServerItem, len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

func (d *SimpleDiscovery) Register(server *registry.ServerItem) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = append(d.servers, server)
}
