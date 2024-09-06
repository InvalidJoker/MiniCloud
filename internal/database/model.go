package database

import (
	"fmt"
	"minicloud/internal/utils"
	"net"

	"go.minekube.com/gate/pkg/edition/java/proxy"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
	Port int    `gorm:"unique"`

	Lobby      bool `gorm:"default:false"`
	Persistent bool `gorm:"default:true"`

	CustomData datatypes.JSON `gorm:"type:json"`

	// template can be nil
	TemplateID string
	Template   Template `gorm:"foreignkey:TemplateID"`
}

type Template struct {
	gorm.Model
	Name     string `gorm:"primaryKey"`
	Software string `gorm:"default:paper"`
	Version  string `gorm:"default:latest"`

	CustomImage     string
	CustomImageData datatypes.JSON `gorm:"type:json"`
}

type CreateServerRequest struct {
	Name string `validate:"required"`
	Port int    `validate:"required"`

	Lobby    bool   // default false
	Template string // can be nil

	// can be nil
	CustomData datatypes.JSON `json:"custom_data"`
}

func (server *Server) GetServerInfo() proxy.ServerInfo {
	ip, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d",
		"0.0.0.0", server.Port))
	if err != nil {
		return nil
	}

	return proxy.NewServerInfo(server.Name, ip)
}

func (t *Template) MoveToServer(server string) error {
	srcDir := fmt.Sprintf("data/templates/%s", t.Name)
	dstDir := fmt.Sprintf("data/servers/%s", server)

	err := utils.CopyDir(srcDir, dstDir)

	if err != nil {
		return err
	}

	return nil

}

func (server *Server) ToRequest() *CreateServerRequest {
	return &CreateServerRequest{
		Name:       server.Name,
		Port:       server.Port,
		Lobby:      server.Lobby,
		Template:   server.Template.Name,
		CustomData: server.CustomData,
	}
}
