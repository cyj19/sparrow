/**
 * @Author: cyj19
 * @Date: 2022/2/28 15:13
 */

package client

import (
	"errors"
	"fmt"
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/protocol"
	"net"
	"sync"
)

type Client struct {
	Option *Option
	rw     *sync.RWMutex
}

func NewClient() *Client {
	return &Client{rw: new(sync.RWMutex)}
}

func (c *Client) Call(serviceName, serviceMethod string, args, reply interface{}) error {

	conn, err := net.Dial("tcp", "0.0.0.0:8787")
	if err != nil {
		return err
	}
	defer conn.Close()

	if serviceName == "" || serviceMethod == "" {
		return errors.New("serviceName or serviceMethod is null")
	}
	// 构建请求
	reqHeader := &protocol.Header{
		Start:     protocol.StartChar,
		Version:   byte(1),
		CodecType: byte(codec.JSON),
	}
	reqBody := &protocol.Body{
		Magic:         "11",
		ServiceName:   serviceName,
		ServiceMethod: serviceMethod,
	}
	// 序列化
	codecPlugin, ok := codec.Get(codec.JSON)
	if !ok {
		return errors.New("codec plugin is not exist")
	}
	metaData, err := codecPlugin.Encode(args)
	if err != nil {
		return errors.New(fmt.Sprintf("client encode metaData error:%v", err))
	}
	reqBody.MetaData = metaData
	payload, err := codecPlugin.Encode(reply)
	if err != nil {
		return errors.New(fmt.Sprintf("client encode payload error:%v", err))
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

	//c.rw.RLock()
	_, err = conn.Write(reqData)
	//c.rw.Unlock()
	return err
}
