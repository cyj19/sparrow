/**
 * @Author: cyj19
 * @Date: 2022/2/24 14:49
 */

package server

import (
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/protocol"
	"log"
	"net"
	"reflect"
)

// 处理请求
func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	for {
		message, err := protocol.DecodeMessage(conn)
		if err != nil {
			log.Printf("protocol.DecodeMessage error:%v", err)
			return
		}
		// 反序列化
		cType := codec.CodecType(message.Header.CodecType)
		codecPlugin, ok := codec.Get(cType)
		if !ok {
			log.Println("rpc not have this codecType")
			return
		}
		var serviceName, serviceMethod string
		err = codecPlugin.Decode([]byte(message.Body.ServiceName), serviceName)
		if err != nil {
			log.Printf("codecPlugin.Decode error:%v", err)
			return
		}
		err = codecPlugin.Decode([]byte(message.Body.ServiceMethod), serviceMethod)
		if err != nil {
			log.Printf("codecPlugin.Decode error:%v", err)
			return
		}
		metaData := make([]byte, len(message.Body.MetaData))
		err = codecPlugin.Decode(message.Body.MetaData, metaData)
		if err != nil {
			log.Printf("codecPlugin.Decode error:%v", err)
			return
		}
		payload := make([]byte, len(message.Body.Payload))
		err = codecPlugin.Decode(message.Body.Payload, payload)
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
			Header: message.Header,
		}
		body := &protocol.Body{
			ServiceName:   message.Body.ServiceName,
			ServiceMethod: message.Body.ServiceMethod,
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
		conn.Write(msgData)
	}

}
