/**
 * @Author: cyj19
 * @Date: 2022/3/4 14:29
 */

package transport

import (
	"net"
	"time"
)

func init() {
	Server.register(TCP, defaultTcpServer)
	Client.register(TCP, defaultTcpClient)
}

func defaultTcpServer(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func defaultTcpClient(addr string, timeout time.Duration) (net.Conn, error) {
	if timeout > 0 {
		return net.DialTimeout("tcp", addr, timeout)
	}
	return net.Dial("tcp", addr)
}
