package events

import (
	"minicloud/internal/cloud"
	"minicloud/internal/database"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type EventHandlers struct {
	Database *database.Database
	Proxy    *proxy.Proxy
	Docker   *cloud.DockerService
}

func NewEventHandlers(db *database.Database, p *proxy.Proxy, d *cloud.DockerService) *EventHandlers {
	return &EventHandlers{
		Database: db,
		Proxy:    p,
		Docker:   d,
	}
}
