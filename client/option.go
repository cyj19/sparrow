/**
 * @Author: cyj19
 * @Date: 2022/2/28 15:27
 */

package client

import (
	"github.com/cyj19/sparrow/balance"
	"github.com/cyj19/sparrow/codec"
	"github.com/cyj19/sparrow/compressor"
	"time"
)

// Option 客户端配置
type Option struct {
	loadBalance    balance.LoadBalancing     // 负载均衡插件
	codecType      codec.CodecType           // 序列化插件
	compressorType compressor.CompressorType // 压缩插件
	readTimeout    time.Duration             // io读取超时时间
	writeTimeout   time.Duration             // io写超时时间
	connectTimeout time.Duration             // 连接超时时间
}

func defaultOption() *Option {
	return &Option{
		loadBalance:    balance.NewRoundRobin(),
		codecType:      codec.JSON,
		compressorType: compressor.GZIP,
		readTimeout:    3 * time.Minute,
		writeTimeout:   1 * time.Minute,
		connectTimeout: 1 * time.Minute,
	}
}
