/**
 * @Author: cyj19
 * @Date: 2022/2/28 16:15
 */

package main

import (
	"context"
	"fmt"
	"github.com/cyj19/sparrow/balance"
	"github.com/cyj19/sparrow/client"
	"github.com/cyj19/sparrow/discovery"
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
	d := discovery.NewSparrowDiscovery("http://localhost:9999/sparrow/registry", 0, balance.NewRoundRobin())
	c, err := client.NewClient(d)
	if err != nil {
		log.Fatalln(err)
	}
	wg := sync.WaitGroup{}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)

	for i := 1; i < 11; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("cyj%d", i)
			reqArgs := RequestArg{Name: name}
			respReply := ResponseReply{}
			err := c.Call(ctx, "HelloWorld", "Hello", &reqArgs, &respReply)
			if err != nil {
				log.Printf("call error:%v", err)
			}
			fmt.Println(respReply)
		}(i)
	}

	wg.Wait()

}
