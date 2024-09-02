package events

import (
	"fmt"
	"minicloud/internal/database"

	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type EventHandlers struct {
	Database *database.Database
	Proxy    *proxy.Proxy
}

func NewEventHandlers(db *database.Database, p *proxy.Proxy) *EventHandlers {
	return &EventHandlers{
		Database: db,
		Proxy:    p,
	}
}

var serversMap = make(map[proxy.Player]proxy.RegisteredServer)

func (e *EventHandlers) HandlePlayerJoin(event *proxy.PlayerChooseInitialServerEvent) {

	// Check if the player has already been assigned a server for performance reasons
	if server, ok := serversMap[event.Player()]; ok {
		event.SetInitialServer(server)
		return
	}

	var server *database.Server
	e.Database.Where(&server, "lobby = ?", true)

	if server == nil {
		msg := "There are currently no available servers. Please try again later."
		player := event.Player()
		player.Disconnect(&component.Text{Content: msg})
	}

	var regServer proxy.RegisteredServer

	servers := e.Proxy.Servers()
	for _, gateServer := range servers {
		if gateServer.ServerInfo().Name() == server.Name {
			regServer = gateServer
			break
		}
	}

	fmt.Println("Server: ", server.Name)

	if regServer != nil {
		event.SetInitialServer(regServer)
		serversMap[event.Player()] = regServer
	} else {
		msg := "There are currently no available servers. Please try again later."
		player := event.Player()
		player.Disconnect(&component.Text{Content: msg})
	}

}
