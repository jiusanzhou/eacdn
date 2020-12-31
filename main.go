package main

import (
	"log"

	"go.zoe.im/x/cli"

	"go.zoe.im/eacdn/cmd"
	"go.zoe.im/eacdn/pkg/service"

	_ "go.zoe.im/eacdn/pkg/provider/caddy"
	_ "go.zoe.im/eacdn/pkg/provider/nginx"
)

func main() {
	svr := service.New()
	cmd.Option(
		cli.GlobalConfig(svr.Config, cli.WithConfigName()),
		cli.Run(func(c *cli.Command, args ...string) {
			if err := svr.Run(); err != nil {
				log.Println("start the eacdn service error:", err)
			}
		}),
	)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
