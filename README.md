# sparrow
sparrow意为麻雀，俗话说“麻雀虽小，五脏俱全”，希望该项目也如此，能帮助大家更好地理解RPC原理

## 上手指南
以下指南帮助你在本地机器上运行该项目

### 下载项目
```
git clone https://github.com/cyj19/sparrow.git
```

### 运行项目
```
# 进入示例
cd examples/helloworld

# 启动注册中心
cd registry
go run registry.go

# 启动服务端
cd server
go run server.go

# 启动客户端
cd client
go run client.go
```

### 如何使用sparrow
1. 启动注册中心registry  
2. 在服务端引用sparrow的server 
```
import (
    "github.com/cyj19/sparrow/server"
    "github.com/cyj19/sparrow/registry"
)

func main() {
    // 创建rpc服务端
	s := server.NewServer()
	// 注册服务
	s.Register(&HelloWorld{})
	// 向注册中心发送心跳
	registry.HeartBeat("http://localhost:9999/sparrow/registry", "tcp", ":8787", 0)
	// 启动rpc服务端
	err := s.Run(server.UseTCP("0.0.0.0:8787"))
	if err != nil {
		log.Fatalln(err)
	}
}
```

3. 在客户端引用sparrow的client  
```
import (
    "github.com/cyj19/saprrow/client"
    "github.com/cyj19/sparrow/balance"
    "github.com/cyj19/sparrow/discovery"
    )

func main() {
    // discovery用于服务发现
	d := discovery.NewSparrowDiscovery("http://localhost:9999/sparrow/registry", 0, balance.NewRoundRobin())
	// 创建rpc客户端
	c, err := client.NewClient(d)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)

	reqArgs := RequestArg{Name: "cyj19"}
	respReply := ResponseReply{}
	// 远程调用
	err = c.Call(ctx, "HelloWorld", "Hello", &reqArgs, &respReply)
	if err != nil {
		log.Printf("call error:%v", err)
	}
	fmt.Println(respReply)

}

```
