/**
 * @Author: cyj19
 * @Date: 2022/3/1 17:12
 */

package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Gzip struct {
}

func (g *Gzip) Zip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(data)
	if err != nil {
		return nil, err
	}
	err = gw.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *Gzip) Unzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.Write(data)
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(&buf)
	if err != nil {
		return nil, err
	}

	result := bytes.NewBuffer(nil)
	_, err = io.Copy(result, gr)
	if err != nil {
		return nil, err
	}
	err = gr.Close()
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
