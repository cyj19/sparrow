/**
 * @Author: cyj19
 * @Date: 2022/2/24 14:45
 */

package server

import (
	"errors"
	"fmt"
	"github.com/cyj19/sparrow/transport"
	"log"
)

// Server 服务管理器
type Server struct {
	serviceMap map[string]*service // 服务注册
	Option     *Option             // 管理器配置
}

func NewServer() *Server {
	return &Server{
		serviceMap: map[string]*service{},
		Option:     genDefaultOption(),
	}
}

func (s *Server) register(v interface{}, serviceName string, useName bool) error {
	srv, err := newService(v, serviceName, useName)
	if err != nil {
		return err
	}
	if _, ok := s.serviceMap[srv.name]; ok {
		return errors.New(fmt.Sprintf("the service:%s is registered", srv.name))
	}
	s.serviceMap[srv.name] = srv

	return nil
}

func (s *Server) Register(v interface{}) error {
	return s.register(v, "", false)
}

func (s Server) RegisterName(v interface{}, serviceName string) error {
	return s.register(v, serviceName, true)
}

func (s *Server) run() error {
	for {
		// 等待连接
		conn, err := s.Option.nl.Accept()
		if err != nil {
			log.Printf("lis.Accept error:%#v \n", err)
			continue
		}

		// 处理请求
		go s.process(conn)
	}

	return nil
}

func (s *Server) Run(fns ...OptionSetter) error {
	for _, fn := range fns {
		fn(s.Option)
	}
	var err error
	s.Option.nl, err = transport.Server.Gen(s.Option.Protocol, s.Option.Host)
	if err != nil {
		return err
	}

	return s.run()
}
