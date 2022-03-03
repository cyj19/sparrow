/**
 * @Author: cyj19
 * @Date: 2022/2/28 16:15
 */

package main

import (
	"context"
	"fmt"
	"github.com/cyj19/sparrow/client"
	"log"
	"sync"
	"time"
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
	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)

	go func() {
		defer wg.Done()
		reqArgs := RequestArg{Name: "cyj"}
		respReply := ResponseReply{}
		err := c.Call(ctx, "HelloWorld", "Hello", &reqArgs, &respReply)
		if err != nil {
			log.Printf("call error:%v", err)
		}
		fmt.Println(respReply)
	}()

	go func() {
		defer wg.Done()
		reqArgs := RequestArg{Name: "zhangsan"}
		respReply := ResponseReply{}
		err := c.Call(ctx, "HelloWorld", "Hello", &reqArgs, &respReply)
		if err != nil {
			log.Printf("call error:%v", err)
		}
		fmt.Println(respReply)
	}()

	wg.Wait()

}
