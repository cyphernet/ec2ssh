package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetService(profile string, region string) *ec2.EC2 {
	creds := *credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{Filename: "", Profile: profile},
		})

	sessionAws, _ := session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	})

	_, err := creds.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	svc := ec2.New(sessionAws, &aws.Config{Region: aws.String(region), Credentials: &creds})

	return svc
}
