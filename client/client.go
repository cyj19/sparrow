/**
 * @Author: cyj19
 * @Date: 2022/2/28 15:13
 */

package client

import (
	"errors"
	"fmt"
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/compressor"
	"github.com/cyj19/sparrow/protocol"
	"log"
	"net"
	"sync"
)

type Caller struct {
	Reply interface{} // 调用结果
	done  chan error  // 通知调用结束
}

type Client struct {
	Option    *Option
	reqMutex  *sync.Mutex
	respMutex *sync.Mutex
	conn      net.Conn
	callMap   map[string]*Caller
	close     chan error // 通知关闭连接
}

func NewClient(proto, addr string) (*Client, error) {
	conn, err := net.Dial(proto, addr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		reqMutex:  new(sync.Mutex),
		respMutex: new(sync.Mutex),
		conn:      conn,
		callMap:   map[string]*Caller{},
		close:     make(chan error),
	}
	go c.receive()
	return c, nil
}

func (c *Client) Call(serviceName, serviceMethod string, args, reply interface{}) error {
	var err error
	defer func() {
		if err != nil {
			c.conn.Close()
		}
	}()
	if serviceName == "" || serviceMethod == "" {
		return errors.New("serviceName or serviceMethod is null")
	}

	done := make(chan error, 0)
	c.call(done, serviceName, serviceMethod, args, reply)
	select {
	case err, ok := <-done:
		if !ok {
			return nil
		}
		return err
	case <-c.close:
		return errors.New("connect close")
	}

	return nil
}

func (c *Client) call(done chan error, serviceName, serviceMethod string, args, reply interface{}) error {
	// 构建请求
	reqHeader := &protocol.Header{
		Start:     protocol.StartChar,
		Version:   byte(1),
		CodecType: byte(codec.JSON),
	}
	reqBody := &protocol.Body{
		ServiceName:   serviceName,
		ServiceMethod: serviceMethod,
	}

	c.respMutex.Lock()
	c.callMap[reqBody.Magic] = &Caller{
		Reply: reply,
		done:  done,
	}
	c.respMutex.Unlock()

	// 序列化
	codecPlugin, ok := codec.Get(codec.JSON)
	if !ok {
		return errors.New("codec plugin is not exist")
	}
	payload, err := codecPlugin.Encode(args)
	if err != nil {
		return errors.New(fmt.Sprintf("client encode payload error:%v", err))
	}
	// 压缩
	cpr, ex := compressor.Get(compressor.GZIP)
	if !ex {
		return errors.New("compress plugin is not exist")
	}
	payload, err = cpr.Zip(payload)
	if err != nil {
		return errors.New(fmt.Sprintf("client compress payload error:%#v", err))
	}
	reqBody.Payload = payload
	reqMsg := &protocol.Message{
		Header: reqHeader,
		Body:   reqBody,
	}

	reqData, err := protocol.EncodeMessage(reqMsg)
	if err != nil {
		return errors.New(fmt.Sprintf("client encode message error:%v", err))
	}

	c.reqMutex.Lock()
	_, err = c.conn.Write(reqData)
	c.reqMutex.Unlock()

	return err
}

func (c *Client) receive() {
	defer c.conn.Close()
	for {
		magic, callDone, err := c.handleResponse()
		if err != nil && callDone != nil {
			callDone <- err
		}
		if err == nil && magic != "" {

		}

	}
}

func (c *Client) handleResponse() (magic string, done chan error, err error) {
	msg, err := protocol.DecodeMessage(c.conn)
	if err != nil {
		close(c.close)
		return "", nil, err
	}
	caller, _ := c.callMap[msg.Body.Magic]
	// 反序列化
	cType := codec.CodecType(msg.Header.CodecType)
	codecPlugin, ok := codec.Get(cType)
	if !ok {
		err = errors.New("codec plugin is not exist")
		caller.done <- err
		return "", nil, err
	}
	err = codecPlugin.Decode(msg.Body.Payload, caller.Reply)
	if err != nil {
		log.Printf("client decode error:%#v", err)
		return "", nil, err
	}
	return msg.Body.Magic, caller.done, nil
}
