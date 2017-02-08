package cmd

import (
	"fmt"
	"os/exec"
	"os"

	"github.com/cyphernet/ec2ssh/config"
)

func Exec(name, ip string) {
	fmt.Println("Connecting to " + name + " (" + ip + ")")

	username, key := config.New()

	command := exec.Command("ssh", "-i", key, username+"@"+ip)
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	err := command.Run()
	if err != nil {
		panic(err)
	}
}
