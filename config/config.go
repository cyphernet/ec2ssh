package config

import (
	"os"
	"io/ioutil"
	"os/user"
	"fmt"

	"github.com/vaughan0/go-ini"
)

var username = ""
var key = ""

func New() (string, string) {
	usr, _ := user.Current()
	filePath := usr.HomeDir + string(os.PathSeparator) + ".ec2ssh"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		createConfigFile(filePath)
	}

	file, err := ini.LoadFile(filePath)
	if err != nil {
		fmt.Println("Could not load config file " + filePath)
		os.Exit(0)
	}

	username, _ = file.Get("ssh", "username")
	key, _ = file.Get("ssh", "key")

	return username, key
}

func createConfigFile(path string) {
	defaultConfig := "[ssh]\nusername = \nkey ="
	err := ioutil.WriteFile(path, []byte(defaultConfig), 0644)
	if err != nil {
		fmt.Println("Could not create config file " + path)
		os.Exit(0)
	}
}
