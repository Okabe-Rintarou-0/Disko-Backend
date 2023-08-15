package main

import (
	"disko/internal/config"
	"disko/internal/handler"
	"disko/internal/svc"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/disko-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	fmt.Printf("files will be save in directory: %s\n", c.FileStorage.LocalStorage.Dir)

	err := os.Mkdir(c.FileStorage.LocalStorage.Dir, 0700)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(err)
	}

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, nil, "http://localhost:3000"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
