package cloud

import (
	"context"
	"minicloud/internal/database"
)

func (s *DockerService) RegisterServer(ctx context.Context, server *database.Server) error {
	_, err := s.Proxy.Register(server.GetServerInfo())

	if err != nil {
		return err
	}

	return nil

}
