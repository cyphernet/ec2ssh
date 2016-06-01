package main

import (
	"sort"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type Instance struct {
	Name         string
	IP           string
	ID           string
	InstanceType string
	State        string
}

type Instances []Instance

func (slice Instances) Len() int {
	return len(slice)
}

func (slice Instances) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice Instances) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func getInstances(svc *ec2.EC2, all bool) Instances {

	instances := Instances{}

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	for idx := range resp.Reservations {
		for res, inst := range resp.Reservations[idx].Instances {
			name := ""
			for _, tag := range resp.Reservations[idx].Instances[res].Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}
			if inst.PrivateIPAddress != nil {
				if all || *inst.State.Name == "running" {
					instances = append(instances, Instance{name, *inst.PrivateIPAddress, *inst.InstanceID, *inst.InstanceType, *inst.State.Name})
				}
			}
		}
	}

	sort.Sort(instances)

	return instances
}
