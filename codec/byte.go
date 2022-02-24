/**
 * @Author: cyj19
 * @Date: 2022/2/24 16:13
 */

package codec

import (
	"errors"
	"fmt"
	"reflect"
)

type ByteCodec struct {
}

func (b *ByteCodec) Encode(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}
	if data, ok := v.(*[]byte); ok {
		return *data, nil
	}

	return nil, errors.New(fmt.Sprintf("%T cannot encode to []byte", v))
}

func (b *ByteCodec) Decode(data []byte, v interface{}) error {
	reflect.Indirect(reflect.ValueOf(v)).SetBytes(data)
	return nil
}
