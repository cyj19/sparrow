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
	_, err := io.ReadFull(r, headerData)
	if err != nil {
		return nil, err
	}
	// 判断是否属于自定义的消息
	if headerData[0] != StartChar {
		return nil, errors.New("the message is not valid")
	}
	header, err := DecodeHeader(headerData)
	if err != nil {
		return nil, err
	}
	bodySize := header.MagicSize + header.ServiceNameSize + header.ServiceMethodSize + header.MetaDataSize + header.PayLoadSize
	bodyData := make([]byte, bodySize)
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

	body.Magic = string(data[:magicSize])
	body.ServiceName = string(data[magicSize:serviceNameSize])
	body.ServiceMethod = string(data[serviceNameSize:serviceMethodSize])
	body.MetaData = data[serviceMethodSize:metaDataSize]
	body.Payload = data[metaDataSize:payloadSize]
	return body, nil
}

// EncodeMessage 发送前编码消息
func EncodeMessage(message *Message) ([]byte, error) {
	data := make([]byte, 0)
	headerData, err := EncodeHeader(message.Header)
	if err != nil {
		return nil, err
	}
	bodyData, err := EncodeBody(message.Body, message.Header)
	if err != nil {
		return nil, err
	}
	copy(data, headerData)
	copy(data, bodyData)
	return data, nil
}

func EncodeHeader(header *Header) ([]byte, error) {
	data := make([]byte, HeaderSize)
	data[0] = header.Start
	data[1] = header.Version
	data[2] = header.CodecType
	binary.BigEndian.PutUint32(data[2:6], header.MagicSize)
	binary.BigEndian.PutUint32(data[6:10], header.ServiceNameSize)
	binary.BigEndian.PutUint32(data[10:14], header.ServiceMethodSize)
	binary.BigEndian.PutUint32(data[14:18], header.MetaDataSize)
	binary.BigEndian.PutUint32(data[18:22], header.PayLoadSize)
	return data, nil
}

func EncodeBody(body *Body, header *Header) ([]byte, error) {
	bodySize := header.MagicSize + header.ServiceNameSize + header.ServiceMethodSize + header.MetaDataSize + header.PayLoadSize
	data := make([]byte, bodySize)
	data = append(data, body.Magic...)
	data = append(data, body.ServiceName...)
	data = append(data, body.ServiceMethod...)
	data = append(data, body.MetaData...)
	data = append(data, body.Payload...)
	return data, nil
}
