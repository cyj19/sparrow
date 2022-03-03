/**
 * @Author: cyj19
 * @Date: 2022/2/24 14:49
 */

package server

import (
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/compressor"
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
	defer func() {
		sChannel.Close()
	}()

	// 回复消息
	go func() {
	loop:
		for {
			select {
			case respMsg, ok := <-sChannel.Ch:
				if !ok {
					break loop

				}
				// 写入响应
				_, err := conn.Write(respMsg)
				if err != nil {
					log.Printf("conn.Write error:%v", err)
					break loop
				}
				log.Println("finish write message...")
			}
		}
	}()

	// 读取消息
	for {
		message, err := protocol.DecodeMessage(conn)
		if err != nil {
			// 说明连接被对端关闭了
			if err == io.EOF {
				log.Printf("ip: %s close", conn.RemoteAddr())
				break
			}
			//log.Printf("protocol.DecodeMessage error:%v", err)
			break
		}

		go s.handleRequest(sChannel, message)
	}

}

func (s *Server) handleRequest(sChannel *SendChannel, reqMsg *protocol.Message) {

	compressorType := compressor.CompressorType(reqMsg.Header.CompressorType)
	compressPlugin, ex := compressor.Get(compressorType)
	if !ex {
		log.Println("rpc not have this compressor type")
		return
	}

	cType := codec.CodecType(reqMsg.Header.CodecType)
	codecPlugin, ok := codec.Get(cType)
	if !ok {
		log.Println("rpc not have this codecType")
		return
	}
	serviceName := reqMsg.Body.ServiceName
	serviceMethod := reqMsg.Body.ServiceMethod

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
	argVal := reflect.New(method.argType.Elem()).Interface()
	replyVal := reflect.New(method.replyType.Elem()).Interface()

	// 解压
	var err error
	reqMsg.Body.Payload, err = compressPlugin.Unzip(reqMsg.Body.Payload)
	if err != nil {
		log.Printf("server compressor.Unzip error:%#v", err)
	}

	// 反序列化
	err = codecPlugin.Decode(reqMsg.Body.Payload, argVal)
	if err != nil {
		log.Printf("server codecPlugin.Decode error:%v", err)
		return
	}
	// 调用方法
	reflectValues := method.method.Func.Call([]reflect.Value{srv.refVal, reflect.ValueOf(argVal), reflect.ValueOf(replyVal)})
	errorVal := reflectValues[0].Interface()
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
		Magic:         reqMsg.Body.Magic,
		ServiceName:   serviceName,
		ServiceMethod: serviceMethod,
	}
	// 序列化
	body.Payload, err = codecPlugin.Encode(replyVal)
	if err != nil {
		log.Printf("codecPlugin.Encode error:%v", err)
		return
	}
	// 压缩
	body.Payload, err = compressPlugin.Zip(body.Payload)
	if err != nil {
		log.Printf("compressPlugin.Encode error:%v", err)
		return
	}
	msg.Body = body
	msgData, err := protocol.EncodeMessage(msg)
	if err != nil {
		log.Printf("protocol.EncodeMessage error:%v", err)
		return
	}
	err = sChannel.Send(msgData)
	if err != nil {
		log.Println(err)
		return
	}
}
