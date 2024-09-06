package cloud

import (
	"fmt"
	"os"
)

func CreateFolder(path string) {
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return
	}
}

func CreateDataFolder() {
	CreateFolder("data")
	CreateFolder("data/templates")
	CreateFolder("data/servers")
	CreateFolder("data/config")
}

func CreateTemplate(name string) string {
	template := fmt.Sprintf("data/templates/%s", name)
	// create template folder
	err := os.MkdirAll(template, os.ModePerm)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return template

}

func CreateServer(name string) (string, error) {
	server := fmt.Sprintf("data/servers/%s", name)
	err := os.MkdirAll(server, os.ModePerm)

	if err != nil {
		return "", err
	}

	return server, nil
}

func DeleteServer(name string) {
	server := fmt.Sprintf("data/servers/%s", name)
	err := os.RemoveAll(server)

	if err != nil {
		fmt.Println(err)
	}

}
