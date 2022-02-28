/**
 * @Author: cyj19
 * @Date: 2022/2/24 14:49
 */

package server

import (
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/protocol"
	"io"
	"log"
	"net"
	"reflect"
)

// 处理请求
func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	sChannel := NewSendChannel(s.Option.SendChannelSize)
	defer sChannel.Close()
	// 回复消息
	go func() {
		for {
			select {
			case respMsg, ok := <-sChannel.Ch:
				if !ok {
					log.Println("connect is closed...")
					continue

				}
				// 写入响应
				_, err := conn.Write(respMsg)
				if err != nil {
					log.Printf("conn.Write error:%v", err)
					continue
				}
			}
		}
	}()

	// 处理消息
	for {
		message, err := protocol.DecodeMessage(conn)
		if err != nil {
			if err == io.EOF {
				//log.Printf("ip: %s close", conn.RemoteAddr())
				continue
			}
			log.Printf("protocol.DecodeMessage error:%v", err)
			continue
		}
		log.Println(message)
		go s.handleRequest(sChannel, message)
	}

}

func (s *Server) handleRequest(sChannel *SendChannel, reqMsg *protocol.Message) {
	// 反序列化
	cType := codec.CodecType(reqMsg.Header.CodecType)
	codecPlugin, ok := codec.Get(cType)
	if !ok {
		log.Println("rpc not have this codecType")
		return
	}
	serviceName := reqMsg.Body.ServiceName
	serviceMethod := reqMsg.Body.ServiceMethod

	metaData := make([]byte, len(reqMsg.Body.MetaData))
	err := codecPlugin.Decode(reqMsg.Body.MetaData, metaData)
	if err != nil {
		log.Printf("codecPlugin.Decode error:%v", err)
		return
	}
	payload := make([]byte, len(reqMsg.Body.Payload))
	err = codecPlugin.Decode(reqMsg.Body.Payload, &payload)
	if err != nil {
		log.Printf("codecPlugin.Decode error:%v", err)
		return
	}

	// 获取服务实例
	srv, ok := s.serviceMap[serviceName]
	if !ok {
		log.Printf("the service:%s is not register", serviceName)
		return
	}
	method, ok := srv.methodMap[serviceMethod]
	if !ok {
		log.Printf("the method:%s is not register", serviceMethod)
		return
	}
	// 创建参数实例
	argVal := reflect.New(method.argType)
	replyVal := reflect.New(method.replyType)
	// 调用方法
	refVals := method.method.Func.Call([]reflect.Value{srv.refVal, argVal, replyVal})
	errorVal := refVals[0].Interface()
	if errorVal != nil {
		// 调用失败
		log.Printf("%s.%s error:%v", serviceName, serviceMethod, errorVal)
		return
	}
	// 写消息
	msg := &protocol.Message{
		Header: reqMsg.Header,
	}
	body := &protocol.Body{
		ServiceName:   reqMsg.Body.ServiceName,
		ServiceMethod: reqMsg.Body.ServiceMethod,
	}
	// 序列化
	body.MetaData, err = codecPlugin.Encode(argVal.Interface())
	if err != nil {
		log.Printf("codecPlugin.Encode error:%v", err)
		return
	}
	body.Payload, err = codecPlugin.Encode(replyVal.Interface())
	if err != nil {
		log.Printf("codecPlugin.Encode error:%v", err)
		return
	}
	msgData, err := protocol.EncodeMessage(msg)
	if err != nil {
		log.Printf("protocol.EncodeMessage error:%v", err)
		return
	}
	sChannel.Send(msgData)
}
