/**
 * @Author: cyj19
 * @Date: 2022/3/4 14:43
 */

package network

import (
	"net"
	"time"
)

func init() {
	Server.register(UNIX, defaultUnixServer)
	Client.register(UNIX, defaultUnixClient)
}

func defaultUnixServer(addr string) (net.Listener, error) {
	return net.Listen("unix", addr)
}

func defaultUnixClient(addr string, timeout time.Duration) (net.Conn, error) {
	if timeout > 0 {
		return net.DialTimeout("unix", addr, timeout)
	}
	return net.Dial("unix", addr)
}
