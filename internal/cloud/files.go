package cloud

import (
	"fmt"
	"os"
)

func CreateDataFolder() {
	os.MkdirAll("data", os.ModePerm)
	os.MkdirAll("data/config", os.ModePerm)
	os.MkdirAll("data/servers", os.ModePerm)
	os.MkdirAll("data/templates", os.ModePerm)
}

func CreateTemplate(name string) string {
	os.MkdirAll("data/templates", os.ModePerm)
	template := fmt.Sprintf("data/templates/%s", name)
	// create template folder
	os.MkdirAll(template, os.ModePerm)
	return template

}

func CreateServer(name string) string {
	os.MkdirAll("data/servers", os.ModePerm)
	server := fmt.Sprintf("data/servers/%s", name)
	os.MkdirAll(server, os.ModePerm)
	return server
}

func DeleteServer(name string) {
	server := fmt.Sprintf("data/servers/%s", name)
	os.RemoveAll(server)
}
