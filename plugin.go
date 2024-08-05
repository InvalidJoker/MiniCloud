package main

import (
	"context"
	"encoding/json"
	"minicloud/cloud"
	"os"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type Config struct {
	AuthToken string `json:"auth_token"`
	Port      int    `json:"port"`
	Interface string `json:"interface"`
}

var Plugin = proxy.Plugin{
	Name: "MiniCloud",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		data, err := os.ReadFile("config.json")
		if err != nil {
			return err
		}
		var config Config

		err = json.Unmarshal(data, &config)

		if err != nil {
			return err
		}

		dockerService, err := cloud.NewDockerService()

		if err != nil {
			return err
		}

		backendService := cloud.NewBackendService(dockerService)

		go backendService.Start()

		return nil
	},
}
