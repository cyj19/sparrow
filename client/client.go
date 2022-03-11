/**
 * @Author: cyj19
 * @Date: 2022/2/28 15:13
 */

package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/compressor"
	"github.com/cyj19/sparrow/discovery"
	"github.com/cyj19/sparrow/protocol"
	"github.com/cyj19/sparrow/transport"
	"github.com/rs/xid"
	"log"
	"net"
	"sync"
	"time"
)

type Caller struct {
	Reply interface{} // 调用结果
	done  chan error  // 通知调用结束
}

type Client struct {
	Option    *Option
	discovery discovery.Discovery
	reqMutex  *sync.Mutex
	respMutex *sync.Mutex
	conn      net.Conn
	callMap   map[string]*Caller
	close     chan error // 通知关闭连接
}

func NewClient(d discovery.Discovery) (*Client, error) {
	c := &Client{
		Option:    defaultOption(),
		discovery: d,
		reqMutex:  new(sync.Mutex),
		respMutex: new(sync.Mutex),
		callMap:   map[string]*Caller{},
		close:     make(chan error),
	}
	serverItem, err := d.Get()
	if err != nil {
		return nil, err
	}
	conn, err := transport.Client.Gen(transport.Protocol(serverItem.Protocol), serverItem.Addr, c.Option.connectTimeout)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	go c.receive()
	return c, nil
}

func (c *Client) registerCall(magic string, caller *Caller) {
	c.respMutex.Lock()
	c.callMap[magic] = caller
	c.respMutex.Unlock()
}

func (c *Client) removeCall(magic string) {
	c.respMutex.Lock()
	delete(c.callMap, magic)
	c.respMutex.Unlock()
}

func (c *Client) Call(ctx context.Context, serviceName, serviceMethod string, args, reply interface{}) error {

	if serviceName == "" || serviceMethod == "" {
		return errors.New("serviceName or serviceMethod is null")
	}

	done := make(chan error, 0)
	// 生成魔法值
	magic := xid.New().String()
	defer func() {
		c.removeCall(magic)
	}()

	go c.call(done, magic, serviceName, serviceMethod, args, reply)

	select {
	case <-ctx.Done():
		return errors.New("rpc client: call failed: " + ctx.Err().Error())
	case err, ok := <-done:
		if !ok {
			return nil
		}
		return err
	case err := <-c.close:
		return errors.New("connect closed by error: " + err.Error())
	}

}

func (c *Client) call(done chan error, magic, serviceName, serviceMethod string, args, reply interface{}) {
	// 构建请求
	reqHeader := &protocol.Header{
		Start:          protocol.StartChar,
		Version:        byte(1),
		CodecType:      byte(c.Option.codecType),
		CompressorType: byte(c.Option.compressorType),
	}

	reqBody := &protocol.Body{
		Magic:         magic,
		ServiceName:   serviceName,
		ServiceMethod: serviceMethod,
	}

	// 序列化
	codecPlugin, ok := codec.Get(c.Option.codecType)
	if !ok {
		c.close <- errors.New("codec plugin is not exist")
	}
	payload, err := codecPlugin.Encode(args)
	if err != nil {
		c.close <- errors.New(fmt.Sprintf("client encode payload error:%v", err))
	}
	// 压缩
	cpr, ex := compressor.Get(c.Option.compressorType)
	if !ex {
		c.close <- errors.New("compress plugin is not exist")
	}
	payload, err = cpr.Zip(payload)
	if err != nil {
		c.close <- errors.New(fmt.Sprintf("client compress payload error:%#v", err))
	}
	reqBody.Payload = payload
	reqMsg := &protocol.Message{
		Header: reqHeader,
		Body:   reqBody,
	}

	reqData, err := protocol.EncodeMessage(reqMsg)
	if err != nil {
		c.close <- errors.New(fmt.Sprintf("client encode message error:%v", err))
	}

	// 设置写超时
	if c.Option.writeTimeout > 0 {
		now := time.Now()
		_ = c.conn.SetWriteDeadline(now.Add(c.Option.writeTimeout))
	}

	c.reqMutex.Lock()
	_, err = c.conn.Write(reqData)
	c.reqMutex.Unlock()

	if err != nil {
		c.close <- err
	}

	c.registerCall(magic, &Caller{
		Reply: reply,
		done:  done,
	})

}

func (c *Client) receive() {
	defer func() {
		_ = c.conn.Close()
	}()
	for {
		callDone, err := c.handleResponse()
		// 调用出错
		if err != nil && callDone != nil {
			callDone <- err
		}
		// 客户端发生错误
		if err != nil && callDone == nil {
			c.close <- err
			break
		}
		// 正常调用结束
		if err == nil && callDone != nil {
			close(callDone)
		}

	}
}

func (c *Client) handleResponse() (done chan error, err error) {
	// 设置读超时
	if c.Option.readTimeout > 0 {
		now := time.Now()
		_ = c.conn.SetReadDeadline(now.Add(c.Option.readTimeout))
	}
	msg, err := protocol.DecodeMessage(c.conn)
	if err != nil {
		return nil, err
	}
	caller, ex := c.callMap[msg.Body.Magic]
	if !ex {
		err = errors.New("codec plugin is not exist")
		return nil, err
	}
	// 解压
	compressorType := compressor.CompressorType(msg.Header.CompressorType)
	compressPlugin, ex := compressor.Get(compressorType)
	if !ex {
		err = errors.New("compressor plugin is not exist")
		return nil, err
	}
	msg.Body.Payload, err = compressPlugin.Unzip(msg.Body.Payload)
	if err != nil {
		return nil, err
	}
	// 反序列化
	cType := codec.CodecType(msg.Header.CodecType)
	codecPlugin, ok := codec.Get(cType)
	if !ok {
		err = errors.New("codec plugin is not exist")
		return nil, err
	}
	err = codecPlugin.Decode(msg.Body.Payload, caller.Reply)
	if err != nil {
		log.Printf("client decode error:%#v", err)
		return caller.done, err
	}
	return caller.done, nil
}
