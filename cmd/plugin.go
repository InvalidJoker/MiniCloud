package main

import (
	"context"
	"encoding/json"
	"minicloud/internal/cloud"
	"minicloud/internal/config"
	"minicloud/internal/database"
	"minicloud/internal/events"
	"os"

	"github.com/robinbraemer/event"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "MiniCloud",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		cloud.CreateDataFolder()
		data, err := os.ReadFile("data/config/config.json")
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

		dockerService, err := cloud.NewDockerService(db, config, p)

		if err != nil {
			return err
		}

		err = dockerService.LoadServers(ctx)

		if err != nil {
			return err
		}

		err = dockerService.DockerTest(&ctx)

		if err != nil {
			return err
		}

		backendService := cloud.NewBackendService(dockerService)

		eventHandler := events.NewEventHandlers(db, p)

		event.Subscribe(p.Event(), 0, onPlayerChooseInitialServer(eventHandler))

		go backendService.Start()

		return nil
	},
}

func onPlayerChooseInitialServer(ev *events.EventHandlers) func(*proxy.PlayerChooseInitialServerEvent) {
	return func(e *proxy.PlayerChooseInitialServerEvent) {
		ev.HandlePlayerJoin(e)
	}
}

func main() {
	proxy.Plugins = append(proxy.Plugins, Plugin)

	gate.Execute()
}
