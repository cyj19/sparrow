/**
 * @Author: cyj19
 * @Date: 2022/3/1 18:35
 */

package compressor

import (
	"log"
	"testing"
)

func TestGzip(t *testing.T) {
	data := []byte("a simple example")
	g := &Gzip{}
	zData, err := g.Zip(data)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(zData))

	uData, err := g.Unzip(zData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(uData))

}
