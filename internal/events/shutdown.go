package events

import (
	"fmt"
	"minicloud/internal/cloud"
	"minicloud/internal/database"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func (e *EventHandlers) HandleShutdown(event *proxy.ShutdownEvent) {
	fmt.Println("Server shutting down")
	// TODO: delete server from database

	var servers []*database.Server
	e.Database.Find(&servers)

	for _, server := range servers {
		if server.Lobby {
			// only stop lobby servers
			err := e.Docker.StopServer(e.Docker.Context, server)

			if err != nil {
				fmt.Println(err)
			}
		} else {
			// delete servers from docker but not from database!!!
			err := e.Docker.DeleteServer(e.Docker.Context, server)

			if err != nil {
				fmt.Println(err)
			}
		}

		if !server.Persistent {
			e.Database.Delete(server)
			cloud.DeleteServer(server.Name)
		}

	}

}
