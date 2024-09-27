package main

import (
	"flag"
	"fmt"
	"zero-im/apps/im/ws/internal/config"
	"zero-im/apps/im/ws/internal/handler"
	"zero-im/apps/im/ws/internal/svc"
	"zero-im/apps/im/ws/websocket"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	srv := websocket.NewServer(c.ListenOn)
	defer srv.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(srv, ctx)

	fmt.Printf("Starting server at %s...\n", c.ListenOn)
	srv.Start()
}
