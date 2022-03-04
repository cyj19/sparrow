/**
 * @Author: cyj19
 * @Date: 2022/2/24 14:56
 */

package server

import (
	"context"
	"github.com/cyj19/sparrow/network"
	"net"
)

// Option 每个服务端的配置
type Option struct {
	ctx             context.Context
	Protocol        network.Protocol // 通信协议
	Host            string           // 服务端地址
	nl              net.Listener
	SendChannelSize int
}

func genDefaultOption() *Option {
	return &Option{
		ctx:             context.Background(),
		Protocol:        network.TCP,
		Host:            "0.0.0.0:8787",
		SendChannelSize: 1000,
	}
}

// OptionSetter 快速设置Option
type OptionSetter func(option *Option)

func UseTCP(host string) OptionSetter {
	return func(option *Option) {
		option.Host = host
		option.Protocol = network.TCP
	}
}

func UseUnix(host string) OptionSetter {
	return func(option *Option) {
		option.Host = host
		option.Protocol = network.UNIX
	}
}

func UseHTTP(host string) OptionSetter {
	return func(option *Option) {
		option.Host = host
		option.Protocol = "http"
	}
}
