/**
 * @Author: cyj19
 * @Date: 2022/2/24 15:54
 */

package codec

import "encoding/json"

type JsonCodec struct {
}

var _ Codec = (*JsonCodec)(nil)

func (j *JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j JsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
