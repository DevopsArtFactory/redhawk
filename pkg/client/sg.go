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
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
)

type SGClient struct {
	Resource string
	Client   *ec2.EC2
}

// GetResourceName returns resource name of client
func (s SGClient) GetResourceName() string {
	return s.Resource
}

// NewSGClient creates a SGClient
func NewSGClient(helper Helper) (Client, error) {
	session := GetAwsSession()
	return &SGClient{
		Resource: constants.SGResourceName,
		Client:   GetEC2ClientFn(session, helper.Region, helper.Credentials),
	}, nil
}

// Scan scans all data
func (s SGClient) Scan() ([]resource.Resource, error) {
	var wg sync.WaitGroup
	var result []resource.Resource

	securityGroups, err := s.GetSGList()
	if err != nil {
		return nil, err
	}

	if len(securityGroups) == 0 {
		logrus.Debug("no security group found")
		return nil, nil
	}

	input := make(chan resource.SGResource)
	output := make(chan []resource.Resource)
	defer close(output)

	go func(input chan resource.SGResource, output chan []resource.Resource, wg *sync.WaitGroup) {
		var ret []resource.Resource
		for result := range input {
			ret = append(ret, result)
			wg.Done()
		}

		output <- ret
	}(input, output, &wg)

	f := func(sg *ec2.SecurityGroup, ch chan resource.SGResource) {
		tmp := resource.SGResource{
			ResourceType: aws.String(constants.SGResourceName),
		}

		tmp.Name = sg.GroupName
		tmp.ID = sg.GroupId
		tmp.VpcID = sg.VpcId
		tmp.Owner = sg.OwnerId
		tmp.Description = sg.Description
		inboundCount := 0
		for _, in := range sg.IpPermissions {
			inboundCount += len(in.IpRanges)
			inboundCount += len(in.UserIdGroupPairs)
		}

		outboundCount := 0
		for _, out := range sg.IpPermissionsEgress {
			outboundCount += len(out.IpRanges)
			outboundCount += len(out.UserIdGroupPairs)
		}

		tmp.InboundCount = aws.Int(inboundCount)
		tmp.OutboundCount = aws.Int(outboundCount)

		ch <- tmp
	}

	logrus.Debugf("Security group found: %d", len(securityGroups))
	for _, sg := range securityGroups {
		wg.Add(1)
		go f(sg, input)
	}

	wg.Wait()
	close(input)

	result = <-output
	logrus.Debugf("total valid Security group data count: %d", len(result))

	return result, nil
}

// GetSGList returns all security group list in the account
func (s SGClient) GetSGList() ([]*ec2.SecurityGroup, error) {
	result, err := s.Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return nil, err
	}

	return result.SecurityGroups, nil
}
