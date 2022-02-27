/**
 * @Author: cyj19
 * @Date: 2022/2/24 15:29
 */

package codec

func init() {
	defaultManager.register(JSON, &JsonCodec{})
	defaultManager.register(BYTE, &ByteCodec{})
}

const (
	JSON CodecType = iota
	BYTE
)

// Codec 序列化接口
type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

type CodecType byte

type codecManager struct {
	codecMap map[CodecType]Codec
}

var defaultManager = &codecManager{
	codecMap: make(map[CodecType]Codec),
}

func (m *codecManager) register(cType CodecType, codec Codec) {
	m.codecMap[cType] = codec
}

func Register(cType CodecType, codec Codec) {
	defaultManager.register(cType, codec)
}

func Get(cType CodecType) (Codec, bool) {
	c, ok := defaultManager.codecMap[cType]
	return c, ok
}
