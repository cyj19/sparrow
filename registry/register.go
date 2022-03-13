/**
 * @Author: cyj19
 * @Date: 2022/3/13 16:59
 */

package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type ServerItem struct {
	Protocol string
	Addr     string
	start    time.Time // 注册时间
}

type SparrowRegistry struct {
	mu      *sync.Mutex
	servers map[string]*ServerItem
	timeout time.Duration
}

const (
	defaultRegistryAddr = "/sparrow/registry"
	defaultTimeout      = 5 * time.Minute
)

var DefaultSparrowRegistry = New(defaultTimeout)

func New(timeout time.Duration) *SparrowRegistry {
	return &SparrowRegistry{
		mu:      new(sync.Mutex),
		servers: make(map[string]*ServerItem),
		timeout: timeout,
	}
}

func (r *SparrowRegistry) putServer(protocol, addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s@%s", protocol, addr)
	_, ok := r.servers[key]
	if !ok {
		r.servers[key] = &ServerItem{
			Protocol: protocol,
			Addr:     addr,
			start:    time.Now(),
		}
	} else {
		r.servers[key].start = time.Now() // 更新服务注册时间
	}
}

func (r *SparrowRegistry) aliveServers() []*ServerItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	var servers []*ServerItem
	for key, server := range r.servers {
		// 服务未过期
		if r.timeout == 0 || server.start.Add(r.timeout).After(time.Now()) {
			servers = append(servers, server)
		} else {
			// 删除过期服务
			delete(r.servers, key)
		}
	}
	return servers
}

func (r *SparrowRegistry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		result, err := json.Marshal(r.aliveServers())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(result)
	case http.MethodPost:
		param, err := io.ReadAll(req.Body)
		if err != nil {
			if err != io.EOF {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		var server ServerItem
		err = json.Unmarshal(param, &server)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		r.putServer(server.Protocol, server.Addr)
		w.WriteHeader(http.StatusOK)
	}
}

func (r *SparrowRegistry) HandleHTTP(registryAddr string) {
	http.Handle(registryAddr, r)
	log.Println("sparrow rpc registry addr: ", registryAddr)
}

func HandleHTTP() {
	DefaultSparrowRegistry.HandleHTTP(defaultRegistryAddr)
}

func sendHeartBeat(registry, protocol, addr string) error {
	log.Println("send heart beat to registry ", registry)
	server := &ServerItem{
		Protocol: protocol,
		Addr:     addr,
	}
	param, err := json.Marshal(server)
	if err != nil {
		return err
	}
	body := bytes.NewReader(param)
	if _, err = http.Post(registry, "application/json;charset=utf-8", body); err != nil {
		log.Println("send heart beat err", err.Error())
		return err
	}
	return nil
}

func HeartBeat(registry, protocol, addr string, timeout time.Duration) {
	if timeout == 0 {
		timeout = defaultTimeout - time.Minute
	}
	err := sendHeartBeat(registry, protocol, addr)
	go func() {
		ticker := time.NewTicker(timeout)
		for err == nil {
			<-ticker.C
			err = sendHeartBeat(registry, protocol, addr)
		}
	}()
}

func Run(protocol, addr string) error {
	if addr == "" {
		return errors.New("registry addr is null")
	}
	l, _ := net.Listen(protocol, addr)
	HandleHTTP()
	if err := http.Serve(l, nil); err != nil {
		return err
	}
	return nil
}
