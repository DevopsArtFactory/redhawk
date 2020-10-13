/*
Copyright 2020 The redhawk Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
)

type STSClient struct {
	Client *sts.STS
}

// NewSTSClient creates a STSClients
func NewSTSClient(region string) (*STSClient, error) {
	session := GetAwsSession()
	return &STSClient{
		Client: GetSTSClientFn(session, region, nil),
	}, nil
}

// GetSTSClientFn creates a new STS Client
func GetSTSClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *sts.STS {
	if creds == nil {
		return sts.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return sts.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// CheckWhoIam calls get-caller-identity and print the result
func (s STSClient) CheckWhoIam() (*string, error) {
	result, err := s.Client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return result.Arn, nil
}
