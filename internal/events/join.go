package events

import (
	"fmt"
	"minicloud/internal/database"

	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var serversMap = make(map[proxy.Player]proxy.RegisteredServer)

func (e *EventHandlers) HandlePlayerJoin(event *proxy.PlayerChooseInitialServerEvent) {

	// Check if the player has already been assigned a server for performance reasons
	if server, ok := serversMap[event.Player()]; ok {
		event.SetInitialServer(server)
		return
	}

	var server *database.Server
	e.Database.Where("lobby = ?", true).First(&server)

	if server == nil || server.Name == "" {
		msg := "There are currently no available servers. Please try again later."
		player := event.Player()
		player.Disconnect(&component.Text{Content: msg})
	}

	var regServer proxy.RegisteredServer

	fmt.Println("Servers: ", e.Proxy.Servers())

	fmt.Println("Server: ", server.Name)
	for _, gateServer := range e.Proxy.Servers() {
		fmt.Println("GateServer: ", gateServer.ServerInfo().Name())
		if gateServer.ServerInfo().Name() == server.Name {
			regServer = gateServer
			break
		}
	}

	if regServer != nil {
		event.SetInitialServer(regServer)
		serversMap[event.Player()] = regServer
	} else {
		msg := "There seems to be an issue with the server. Please try again later."
		player := event.Player()
		player.Disconnect(&component.Text{Content: msg})
	}

}
