/**
 * @Author: cyj19
 * @Date: 2022/2/28 15:27
 */

package client

import (
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/compressor"
)

// Option 客户端配置
type Option struct {
	codecType      codec.CodecType           // 序列化插件
	compressorType compressor.CompressorType // 压缩插件
}

func defaultOption() *Option {
	return &Option{
		codecType:      codec.JSON,
		compressorType: compressor.GZIP,
	}
}
