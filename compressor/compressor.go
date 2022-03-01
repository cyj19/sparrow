/**
 * @Author: cyj19
 * @Date: 2022/3/1 16:38
 */

package compressor

import "errors"

func init() {
	cManager.Register(GZIP, &Gzip{})
}

// Compressor 压缩解压接口
type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}

// CompressorType 压缩类型
type CompressorType byte

type compressorManager struct {
	compressorMap map[CompressorType]Compressor
}

const (
	GZIP CompressorType = iota
)

var cManager = &compressorManager{
	compressorMap: map[CompressorType]Compressor{},
}

func (m *compressorManager) Get(cType CompressorType) (Compressor, bool) {
	compressor, ok := m.compressorMap[cType]
	return compressor, ok
}

func (m *compressorManager) Register(cType CompressorType, compressor Compressor) error {
	return m.register(cType, compressor)
}

func (m *compressorManager) register(cType CompressorType, compressor Compressor) error {
	if _, ok := m.compressorMap[cType]; ok {
		return errors.New("compressor is registered")
	}
	m.compressorMap[cType] = compressor
	return nil
}
