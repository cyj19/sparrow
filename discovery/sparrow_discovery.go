/**
 * @Author: cyj19
 * @Date: 2022/3/13 18:17
 */

package discovery

import (
	"encoding/json"
	"errors"
	"github.com/cyj19/sparrow/balance"
	"github.com/cyj19/sparrow/registry"
	"io"
	"log"
	"net/http"
	"time"
)

type SparrowDiscovery struct {
	*SimpleDiscovery
	registryAddr string
	timeout      time.Duration
	lastUpdate   time.Time
}

var _ Discovery = (*SparrowDiscovery)(nil)

func NewSparrowDiscovery(registryAddr string, timeout time.Duration, lb balance.LoadBalancing) *SparrowDiscovery {
	return &SparrowDiscovery{
		SimpleDiscovery: NewSimpleDiscovery(lb),
		registryAddr:    registryAddr,
		timeout:         timeout,
	}
}

func (d *SparrowDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}
	resp, err := http.Get(d.registryAddr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var servers []*registry.ServerItem
	err = json.Unmarshal(body, &servers)
	if err != nil {
		return err
	}
	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

func (d *SparrowDiscovery) Update(servers []*registry.ServerItem) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

func (d *SparrowDiscovery) Get() (*registry.ServerItem, error) {
	err := d.Refresh()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	n := len(d.servers)
	if n == 0 {
		return nil, errors.New("servers is null")
	}
	index := d.lb.GetModeResult(n)
	return d.servers[index], nil
}

func (d *SparrowDiscovery) GetAll() ([]*registry.ServerItem, error) {
	err := d.Refresh()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	n := len(d.servers)
	if n == 0 {
		return nil, errors.New("servers is null")
	}
	servers := make([]*registry.ServerItem, n, n)
	copy(servers, d.servers)
	return servers, nil
}
