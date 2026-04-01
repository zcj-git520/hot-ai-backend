package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	//"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

var configFile = flag.String("f", "etc/content-svc.yaml", "the config file")

func main() {
	flag.Parse()

	var c zrpc.RpcServerConf
	if err := conf.Load(*configFile, &c); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v\n", err)
		os.Exit(1)
	}

	// TODO: 创建服务上下文
	// ctx := svc.NewServiceContext(c)

	// TODO: 注册 gRPC 服务
	// server := rest.MustNewServer(c)
	// defer server.Stop()
	// pb.RegisterContentServer(server, handler.NewContentHandler(ctx))

	fmt.Printf("Starting content-svc at %s...\n", c.ListenOn)
	// server.Start()
}
