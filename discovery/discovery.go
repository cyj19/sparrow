/**
 * @Author: cyj19
 * @Date: 2022/3/10 18:51
 */

package discovery

import (
	"errors"
	"github.com/cyj19/sparrow/balance"
	"sync"
)

type Discovery interface {
	Refresh() error                     // 从注册中心更新服务列表
	Update(servers []*ServerItem) error // 手动更新服务列表
	Get() (*ServerItem, error)
	GetAll() ([]*ServerItem, error)
}

type ServerItem struct {
	Protocol string
	Addr     string
}

type SimpleDiscovery struct {
	m       *sync.Mutex
	servers []*ServerItem
	lb      balance.LoadBalancing
}

var _ Discovery = (*SimpleDiscovery)(nil)

func NewSimpleDiscovery(lb balance.LoadBalancing) *SimpleDiscovery {
	return &SimpleDiscovery{
		m:       new(sync.Mutex),
		servers: make([]*ServerItem, 0),
		lb:      lb,
	}
}

func (d *SimpleDiscovery) Refresh() error {
	return nil
}

func (d *SimpleDiscovery) Update(servers []*ServerItem) error {
	d.m.Lock()
	defer d.m.Unlock()
	d.servers = servers
	return nil
}

func (d *SimpleDiscovery) Get() (*ServerItem, error) {
	d.m.Lock()
	defer d.m.Unlock()
	n := len(d.servers)
	if n == 0 {
		return nil, errors.New("rpc discovery: no available servers")
	}
	idx := d.lb.GetModeResult(n)
	s := d.servers[idx]
	return s, nil
}

func (d *SimpleDiscovery) GetAll() ([]*ServerItem, error) {
	d.m.Lock()
	defer d.m.Unlock()
	servers := make([]*ServerItem, len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

func (d *SimpleDiscovery) Register(server *ServerItem) {
	d.m.Lock()
	defer d.m.Unlock()
	d.servers = append(d.servers, server)
}
