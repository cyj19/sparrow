/**
 * @Author: cyj19
 * @Date: 2022/2/28 16:15
 */

package main

import (
	"fmt"
	"github.com/cyj19/sparrow/client"
	"log"
)

type RequestArg struct {
	Name string
}

type ResponseReply struct {
	Msg string
}

func main() {
	c, err := client.NewClient("tcp", "0.0.0.0:8787")
	if err != nil {
		log.Fatalln(err)
	}
	reqArgs := RequestArg{Name: "cyj"}
	respReply := ResponseReply{}
	err = c.Call("HelloWorld", "Hello", &reqArgs, &respReply)
	if err != nil {
		log.Printf("call error:%v", err)
	}
	fmt.Println(respReply)
}
