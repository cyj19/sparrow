/**
 * @Author: cyj19
 * @Date: 2022/2/24 14:56
 */

package server

import (
	"context"
	"net"
)

// Option 每个服务端的配置
type Option struct {
	ctx             context.Context
	Protocol        string // 通信协议
	Host            string // 服务端地址
	nl              net.Listener
	SendChannelSize int
}

// OptionSetter 快速设置Option
type OptionSetter func(option *Option)

func UseTCP(host string) OptionSetter {
	return func(option *Option) {
		option.Host = host
		option.Protocol = "tcp"
	}
}

func UseHTTP(host string) OptionSetter {
	return func(option *Option) {
		option.Host = host
		option.Protocol = "http"
	}
}
