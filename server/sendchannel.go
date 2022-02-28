/**
 * @Author: cyj19
 * @Date: 2022/2/28 14:03
 */

package server

import (
	"errors"
	"sync"
)

type SendChannel struct {
	rw    *sync.RWMutex
	Ch    chan []byte
	close bool
}

func NewSendChannel(size int) *SendChannel {
	return &SendChannel{
		rw: new(sync.RWMutex),
		Ch: make(chan []byte, size),
	}
}

func (c *SendChannel) Send(data []byte) error {
	defer c.rw.Unlock()
	c.rw.RLock()

	if c.close {
		return errors.New("sendChannel is closed")
	}
	c.Ch <- data
	return nil

}

func (c *SendChannel) Close() {
	defer c.rw.Unlock()
	c.rw.Lock()
	if c.close == false {
		c.close = true
		close(c.Ch)
	}
}
