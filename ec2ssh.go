package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
	"github.com/vaughan0/go-ini"
)

type password string

func (p password) Password(user string) (password string, err error) {
	return string(p), nil
}

var wg sync.WaitGroup

func main() {

	app := cli.NewApp()
	app.Name = "ec2ssh"
	app.Usage = "ssh into ec2 instances"
	app.Action = func(c *cli.Context) {

		credentials := *credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{Filename: "", Profile: c.String("profile")},
				&credentials.EC2RoleProvider{},
			})

		svc := ec2.New(&aws.Config{Region: c.String("region"), Credentials: &credentials})

		instances := getInstances(svc, (c.Args().First() == "a"))

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "Index\tName\tIP\tInstance Type\tState")
		for i, c := range instances {
			fmt.Fprintln(w, strconv.Itoa(i)+"\t"+c.Name+"\t"+c.IP+"\t"+c.InstanceType+"\t"+c.State)
		}
		w.Flush()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter index: ")
		text, _ := reader.ReadString('\n')
		i, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			panic(err)
		}

		fmt.Println("Connecting to " + instances[i].Name + " (" + instances[i].IP + ")")

		usr, err := user.Current()
		if err != nil {
			panic(err)
		}

		configFile := usr.HomeDir + string(os.PathSeparator) + ".ec2ssh"
		username := ""
		key := ""
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			defaultConfig := "[ssh]\nusername = \nkey ="
			err = ioutil.WriteFile(configFile, []byte(defaultConfig), 0644)
			if err != nil {
				panic(err)
			}

		} else {
			file, err := ini.LoadFile(configFile)
			if err != nil {
				panic(err)
			}

			username, _ = file.Get("ssh", "username")
			key, _ = file.Get("ssh", "key")
		}

		cmd := exec.Command("ssh", "-i", key, username+"@"+instances[i].IP)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "profile, p",
			Value: "default",
			Usage: "aws profile",
		},
		cli.StringFlag{
			Name:  "region, r",
			Value: "us-east-1",
			Usage: "aws region",
		},
	}

	app.Run(os.Args)
}
