package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
)

func GetSTSClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *sts.STS {
	if creds == nil {
		return sts.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return sts.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}
