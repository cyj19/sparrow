/**
 * @Author: cyj19
 * @Date: 2022/2/28 16:19
 */

package main

import (
	"fmt"
	"github.com/cyj19/sparrow/server"
	"log"
)

type HelloWorld struct {
}

type HelloWordRequest struct {
	Name string
}

type HelloWordResponse struct {
	Msg string
}

func (w *HelloWorld) Hello(args HelloWordRequest, reply *HelloWordResponse) error {
	reply.Msg = fmt.Sprintf("hello %s", args.Name)
	return nil
}

func main() {

	s := server.NewServer()
	s.Register(&HelloWorld{})
	err := s.Run(server.UseTCP("0.0.0.0:8787"))
	if err != nil {
		log.Fatalln(err)
	}
}
