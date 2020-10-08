package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds"
)

func GetRDSClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *rds.RDS {
	if creds == nil {
		return rds.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return rds.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}
