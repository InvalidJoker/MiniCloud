package cloud

import (
	"context"
	"fmt"
	"net"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type Server struct {
	ID      string
	Name    string
	Address string
	Port    int
}

func (s *DockerService) RegisterServer(ctx context.Context, server *Server) error {
	ip, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", server.Address, server.Port))
	if err != nil {
		return err
	}

	_, err = s.Proxy.Register(proxy.NewServerInfo(server.Name, ip))

	return err
}
