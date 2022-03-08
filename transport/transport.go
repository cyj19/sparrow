/**
 * @Author: cyj19
 * @Date: 2022/3/4 14:00
 */

package transport

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Protocol string

const (
	TCP  Protocol = "tcp"
	UNIX Protocol = "UNIX"
)

type genServer func(addr string) (net.Listener, error)

type server struct {
	container map[Protocol]genServer
}

var Server = &server{
	container: map[Protocol]genServer{},
}

func (s *server) register(ptl Protocol, fn genServer) {
	s.container[ptl] = fn
}

func (s *server) Gen(ptl Protocol, addr string) (net.Listener, error) {
	fn, ex := s.container[ptl]
	if !ex {
		return nil, errors.New(fmt.Sprintf("rpc do not support protocol: %s", ptl))
	}
	return fn(addr)
}

type genClient func(addr string, timeout time.Duration) (net.Conn, error)

type client struct {
	container map[Protocol]genClient
}

var Client = &client{
	container: map[Protocol]genClient{},
}

func (c *client) register(ptl Protocol, fn genClient) {
	c.container[ptl] = fn
}

func (c *client) Gen(ptl Protocol, addr string, timeout time.Duration) (net.Conn, error) {
	fn, ex := c.container[ptl]
	if !ex {
		return nil, errors.New(fmt.Sprintf("rpc do not support protocol: %s", ptl))
	}
	return fn(addr, timeout)
}
