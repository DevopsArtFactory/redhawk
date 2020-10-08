package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetS3ClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *s3.S3 {
	if creds == nil {
		return s3.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return s3.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}
