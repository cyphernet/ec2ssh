package main

import (
	"os"

	"github.com/cyphernet/ec2ssh/aws"
	"github.com/cyphernet/ec2ssh/instance"
	"github.com/cyphernet/ec2ssh/ui"
	"github.com/cyphernet/ec2ssh/cmd"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "ec2ssh"
	app.Usage = "ssh into ec2 instances"
	app.Action = func(c *cli.Context) {
		svc := aws.GetService(c.String("profile"), c.String("region"))
		ec2Instances := instance.GetInstances(svc, c.Args().First() == "a")

		ui.Show(ec2Instances)
		i := ui.Get()

		cmd.Exec(ec2Instances[i].Name, ec2Instances[i].IP)
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
