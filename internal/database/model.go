package database

import (
	"fmt"
	"minicloud/internal/utils"
	"net"

	"go.minekube.com/gate/pkg/edition/java/proxy"
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
	Port int    `gorm:"unique"`

	Version string `gorm:"default:latest"`

	Lobby bool `gorm:"default:false"`

	TemplateID string
	Template   Template `gorm:"foreignKey:TemplateID"`
}

type Template struct {
	gorm.Model
	Name     string `gorm:"primaryKey"`
	Software string `gorm:"default:paper"`
}

type CreateServerRequest struct {
	Name string
	Port int

	Version  string
	Lobby    bool
	Template string
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

func (server *Server) ToRequest() *CreateServerRequest {
	return &CreateServerRequest{
		Name:     server.Name,
		Port:     server.Port,
		Version:  server.Version,
		Lobby:    server.Lobby,
		Template: server.Template.Name,
	}
}
