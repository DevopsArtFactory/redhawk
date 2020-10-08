package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
)

func GetIAMClientFn(sess client.ConfigProvider, creds *credentials.Credentials) *iam.IAM {
	if creds == nil {
		return iam.New(sess)
	}
	return iam.New(sess, aws.NewConfig().WithCredentials(creds))
}
