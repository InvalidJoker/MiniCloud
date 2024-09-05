package cloud

import (
	"fmt"
	"os"
)

func CreateDataFolder() {
	err := os.MkdirAll("data", os.ModePerm)
	err = os.MkdirAll("data/config", os.ModePerm)
	err = os.MkdirAll("data/servers", os.ModePerm)
	err = os.MkdirAll("data/templates", os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}
}

func CreateTemplate(name string) string {
	err := os.MkdirAll("data/templates", os.ModePerm)
	template := fmt.Sprintf("data/templates/%s", name)
	// create template folder
	err = os.MkdirAll(template, os.ModePerm)

	if err != nil {
		fmt.Println(err)
		return ""
	}

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
