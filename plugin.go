package main

import (
	"context"
	"minicloud/cloud"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "MiniCloud",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		return nil

		dockerService := cloud.NewDockerService()

		


		return nil
	},

}
