package database

import (
	"fmt"
	"minicloud/internal/utils"
	"net"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type Server struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
	Port int    `gorm:"unique"`

	Version string `gorm:"default:latest"`

	Lobby bool `gorm:"default:false"`

	Template Template
}

type Template struct {
	Name     string
	Software string
}

func (server *Server) GetServerInfo() proxy.ServerInfo {
	ip, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d",
		"0.0.0.0", server.Port))
	if err != nil {
		return nil
	}

	return proxy.NewServerInfo(server.Name, ip)
}

func (t *Template) MoveToServer(server string) {
	srcDir := fmt.Sprintf("data/templates/%s", t.Name)
	dstDir := fmt.Sprintf("data/servers/%s", server)

	err := utils.CopyDir(srcDir, dstDir)

	if err != nil {
		fmt.Println(err)
	}

}
