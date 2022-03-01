/**
 * @Author: cyj19
 * @Date: 2022/2/25 21:55
 */

package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

// 消息协议设计 使用前缀长度法
/**
Header:
| start | version | codecType | compressionType | magicSize | serviceNameSize | serviceMethodSize | payloadSize |
| 0x03  |   0x01  |     1     |        1        |     4     |         4       |          4        |      4      |

Body:
| magic | serviceName | serviceMethod | payload |
|   x   |     x       |       x       |    x    |

*/

const (
	HeaderSize = 20
	StartChar  = byte(3)
)

// Header 定义消息头
type Header struct {
	Start             byte   // 起始符
	Version           byte   // 版本号
	CodecType         byte   // 序列化类型
	CompressionType   byte   // 压缩类型
	MagicSize         uint32 // 魔法值大小
	ServiceNameSize   uint32 // 服务名称大小
	ServiceMethodSize uint32 // 服务方法大小
	PayLoadSize       uint32 // 函数参数大小
}

// Body 定义消息体
type Body struct {
	Magic         string // 魔法值
	ServiceName   string // 服务名称
	ServiceMethod string // 服务方法
	Payload       []byte // 函数参数
}

// Message 定义消息
type Message struct {
	Header *Header
	Body   *Body
}

// DecodeMessage 解码消息
func DecodeMessage(r io.Reader) (*Message, error) {
	headerData := make([]byte, HeaderSize)
	// 读取标志位
	_, err := io.ReadFull(r, headerData[:1])
	if err != nil {
		return nil, err
	}
	// 判断是否属于自定义的消息
	if headerData[0] != StartChar {
		return nil, errors.New("the message is not valid")
	}

	// 读取头部剩下的数据
	_, err = io.ReadFull(r, headerData[1:])
	if err != nil {
		return nil, err
	}

	// 解码头部
	header, err := DecodeHeader(headerData)
	if err != nil {
		return nil, err
	}
	bodySize := header.MagicSize + header.ServiceNameSize + header.ServiceMethodSize + header.PayLoadSize
	bodyData := make([]byte, bodySize)
	// 读取消息体的数据
	_, err = io.ReadFull(r, bodyData)
	if err != nil {
		return nil, err
	}

	// 解码消息体
	body, err := DecodeBody(bodyData, header)
	if err != nil {
		return nil, err
	}
	message := &Message{
		Header: header,
		Body:   body,
	}
	return message, nil
}

func DecodeHeader(data []byte) (*Header, error) {
	header := &Header{
		Start:           data[0],
		Version:         data[1],
		CodecType:       data[2],
		CompressionType: data[3],
	}
	// 大端字符序转为uint32
	header.MagicSize = binary.BigEndian.Uint32(data[4:8])
	header.ServiceNameSize = binary.BigEndian.Uint32(data[8:12])
	header.ServiceMethodSize = binary.BigEndian.Uint32(data[12:16])
	header.PayLoadSize = binary.BigEndian.Uint32(data[16:20])
	return header, nil
}

func DecodeBody(data []byte, header *Header) (*Body, error) {
	body := &Body{}

	magicSize := header.MagicSize
	serviceNameSize := header.ServiceNameSize
	serviceMethodSize := header.ServiceMethodSize
	payloadSize := header.PayLoadSize

	var startIndex uint32 = 0
	endIndex := startIndex + magicSize
	length := endIndex - startIndex
	// 创建新的切片进行拷贝，避免操作同一个底层数组
	magic := make([]byte, length)
	copy(magic, data[startIndex:endIndex])
	body.Magic = string(magic)

	startIndex = endIndex
	endIndex = startIndex + serviceNameSize
	length = endIndex - startIndex
	// 创建新的切片进行拷贝，避免操作同一个底层数组
	serviceName := make([]byte, length)
	copy(serviceName, data[startIndex:endIndex])
	body.ServiceName = string(serviceName)

	startIndex = endIndex
	endIndex = startIndex + serviceMethodSize
	length = endIndex - startIndex
	// 创建新的切片进行拷贝，避免操作同一个底层数组
	serviceMethod := make([]byte, length)
	copy(serviceMethod, data[startIndex:endIndex])
	body.ServiceMethod = string(serviceMethod)

	startIndex = endIndex
	endIndex = startIndex + payloadSize
	length = endIndex - startIndex
	payload := make([]byte, length)
	copy(payload, data[startIndex:endIndex])
	body.Payload = payload
	return body, nil
}

// EncodeMessage 发送前编码消息
func EncodeMessage(message *Message) ([]byte, error) {
	header := message.Header
	body := message.Body
	serviceNameByte := []byte(body.ServiceName)
	serviceMethodByte := []byte(body.ServiceMethod)

	msgSize := HeaderSize + len(body.Magic) + len(serviceNameByte) + len(serviceMethodByte) + len(body.Payload)
	data := make([]byte, msgSize)

	// 构建头部
	data[0] = header.Start
	data[1] = header.Version
	data[2] = header.CodecType
	data[3] = header.CompressionType
	binary.BigEndian.PutUint32(data[4:8], uint32(len(body.Magic)))
	binary.BigEndian.PutUint32(data[8:12], uint32(len(serviceNameByte)))
	binary.BigEndian.PutUint32(data[12:16], uint32(len(serviceMethodByte)))
	binary.BigEndian.PutUint32(data[16:20], uint32(len(body.Payload)))

	// 构建body
	startIndex := HeaderSize
	endIndex := startIndex + len(body.Magic)
	copy(data[startIndex:endIndex], body.Magic)

	startIndex = endIndex
	endIndex = startIndex + len(body.ServiceName)
	copy(data[startIndex:endIndex], body.ServiceName)

	startIndex = endIndex
	endIndex = startIndex + len(body.ServiceMethod)
	copy(data[startIndex:endIndex], body.ServiceMethod)

	startIndex = endIndex
	endIndex = startIndex + len(body.Payload)
	copy(data[startIndex:endIndex], body.Payload)

	return data, nil
}
