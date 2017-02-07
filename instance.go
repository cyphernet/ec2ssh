package main

import (
	"sort"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type instance struct {
	Name         string
	IP           string
	ID           string
	InstanceType string
	State        string
}

type instances []instance

func (slice instances) Len() int {
	return len(slice)
}

func (slice instances) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice instances) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func getInstances(svc *ec2.EC2, all bool) instances {

	loadedInstances := instances{}

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
			if inst.PrivateIpAddress != nil {
				if all || *inst.State.Name == "running" {
					loadedInstances = append(loadedInstances, instance{name, *inst.PrivateIpAddress, *inst.InstanceId, *inst.InstanceType, *inst.State.Name})
				}
			}
		}
	}

	sort.Sort(loadedInstances)

	return loadedInstances
}
