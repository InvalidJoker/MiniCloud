package main

import (
	"context"
	"encoding/json"
	"minicloud/cloud"
	"minicloud/config"
	"minicloud/database"
	"os"

	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "MiniCloud",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		data, err := os.ReadFile("config.json")
		if err != nil {
			return err
		}
		var config config.Config

		err = json.Unmarshal(data, &config)

		if err != nil {
			return err
		}

		db, err := database.NewDatabase()

		if err != nil {
			return err
		}

		dockerService, err := cloud.NewDockerService(db, config)

		if err != nil {
			return err
		}

		backendService := cloud.NewBackendService(dockerService)

		go backendService.Start()

		return nil
	},
}

func main() {
	proxy.Plugins = append(proxy.Plugins, Plugin)

	gate.Execute()
}
