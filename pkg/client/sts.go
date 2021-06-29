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
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSClient struct {
	Client *sts.Client
}

// NewSTSClient creates a STSClients
func NewSTSClient(cfg aws.Config) (*STSClient, error) {
	return &STSClient{
		Client: GetSTSClientFn(cfg),
	}, nil
}

// GetSTSClientFn creates a new STS Client
func GetSTSClientFn(cfg aws.Config) *sts.Client {
	return sts.NewFromConfig(cfg)
}

// CheckWhoIam calls get-caller-identity and print the result
func (s STSClient) CheckWhoIam() (*string, error) {
	result, err := s.Client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return result.Arn, nil
}
