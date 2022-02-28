/**
 * @Author: cyj19
 * @Date: 2022/2/25 21:55
 */

package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
)

// 消息协议设计 使用前缀长度法
/**
Header:
| start | version | codecType | magicSize | serviceNameSize | serviceMethodSize | metaDataSize | payloadSize |
| 0x03  |   0x01  |     1     |     4     |         4       |          4        |       4      |      4      |

Body:
| magic | serviceName | serviceMethod | metaData | payload |
|   x   |     x       |       x       |    x     |     x   |

*/

const (
	HeaderSize = 23
	StartChar  = byte(3)
)

// Header 定义消息头
type Header struct {
	Start             byte
	Version           byte
	CodecType         byte
	MagicSize         uint32
	ServiceNameSize   uint32
	ServiceMethodSize uint32
	MetaDataSize      uint32
	PayLoadSize       uint32
}

// Body 定义消息体
type Body struct {
	Magic         string // 大小不确定，用string
	ServiceName   string
	ServiceMethod string
	MetaData      []byte
	Payload       []byte
}

// Message 定义消息
type Message struct {
	Header *Header
	Body   *Body
}

type Protocol struct {
}

func NewProtocol() *Protocol {
	return &Protocol{}
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
	bodySize := header.MagicSize + header.ServiceNameSize + header.ServiceMethodSize + header.MetaDataSize + header.PayLoadSize
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
		Start:     data[0],
		Version:   data[1],
		CodecType: data[2],
	}
	// 大端字符序转为uint32
	header.MagicSize = binary.BigEndian.Uint32(data[2:6])
	header.ServiceNameSize = binary.BigEndian.Uint32(data[6:10])
	header.ServiceMethodSize = binary.BigEndian.Uint32(data[10:14])
	header.MetaDataSize = binary.BigEndian.Uint32(data[14:18])
	header.PayLoadSize = binary.BigEndian.Uint32(data[18:22])
	return header, nil
}

func DecodeBody(data []byte, header *Header) (*Body, error) {
	body := &Body{}

	magicSize := header.MagicSize
	serviceNameSize := header.ServiceNameSize
	serviceMethodSize := header.ServiceMethodSize
	metaDataSize := header.MetaDataSize
	payloadSize := header.PayLoadSize
	log.Println(magicSize, serviceNameSize, serviceMethodSize, metaDataSize, payloadSize)

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
	endIndex = startIndex + metaDataSize
	length = endIndex - startIndex
	metaData := make([]byte, length)
	copy(metaData, data[startIndex:endIndex])
	body.MetaData = metaData

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

	msgSize := HeaderSize + len(body.Magic) + len(serviceNameByte) + len(serviceMethodByte) + len(body.MetaData) + len(body.Payload)
	data := make([]byte, msgSize)

	// 构建头部
	data[0] = header.Start
	data[1] = header.Version
	data[2] = header.CodecType
	binary.BigEndian.PutUint32(data[2:6], uint32(len(body.Magic)))
	binary.BigEndian.PutUint32(data[6:10], uint32(len(serviceNameByte)))
	binary.BigEndian.PutUint32(data[10:14], uint32(len(serviceMethodByte)))
	binary.BigEndian.PutUint32(data[14:18], uint32(len(body.MetaData)))
	binary.BigEndian.PutUint32(data[18:22], uint32(len(body.Payload)))

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
	endIndex = startIndex + len(body.MetaData)
	copy(data[startIndex:endIndex], body.MetaData)

	startIndex = endIndex
	endIndex = startIndex + len(body.Payload)
	copy(data[startIndex:endIndex], body.Payload)

	return data, nil
}
