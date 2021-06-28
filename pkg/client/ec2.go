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
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
)

type EC2Client struct {
	Resource string
	Client   *ec2.Client
	Region   string
	Alias    *string
}

// GetResourceName returns resource name of client
func (e *EC2Client) GetResourceName() string {
	return e.Resource
}

// Scan scans all data
func (e *EC2Client) Scan() ([]resource.Resource, error) {
	var wg sync.WaitGroup
	var result []resource.Resource

	logrus.Debugf("Start to scan all ec2 instances")
	reservations, err := e.GetEC2Instances(nil, nil)
	if err != nil {
		return nil, err
	}

	regionName, err := GetRegionName(e.Region)
	if err != nil {
		return nil, err
	}

	if len(reservations) == 0 {
		logrus.Debug("no ec2 instance found")
		return nil, nil
	}

	input := make(chan resource.EC2Resource)
	output := make(chan []resource.Resource)
	defer close(output)

	go func(input chan resource.EC2Resource, output chan []resource.Resource, wg *sync.WaitGroup) {
		var ret []resource.Resource
		for result := range input {
			ret = append(ret, result)
			wg.Done()
		}

		output <- ret
	}(input, output, &wg)

	f := func(instance types.Instance, ch chan resource.EC2Resource) {
		tmp := resource.EC2Resource{
			ResourceType: aws.String(constants.EC2ResourceName),
		}
		logrus.Tracef("Gathering information about instance: %s", *instance.InstanceId)
		tmp.InstanceID = instance.InstanceId
		tmp.InstanceStatus = aws.String(string(instance.State.Name))
		tmp.InstanceType = aws.String(string(instance.InstanceType))
		tmp.AvailabilityZone = instance.Placement.AvailabilityZone
		tmp.LaunchTime = instance.LaunchTime
		tmp.ImageID = instance.ImageId
		tmp.VpcID = instance.VpcId
		tmp.SubnetID = instance.SubnetId
		tmp.RegionName = regionName
		tmp.AccountAlias = e.Alias

		if instance.PublicIpAddress != nil {
			tmp.PublicIP = instance.PublicIpAddress
		}

		if instance.KeyName != nil {
			tmp.KeyName = instance.KeyName
		}

		//if instance.IamInstanceProfile != nil {
		//	tmp.IAMInstanceProfile = instance.IamInstanceProfile.Arn
		//}

		for _, tag := range instance.Tags {
			if *tag.Key == "Name" {
				tmp.Name = tag.Value
				break
			}
		}

		var privateIps []string
		//var ipv6s []string
		for _, net := range instance.NetworkInterfaces {
			//for _, v6 := range net.Ipv6Addresses {
			//	ipv6s = append(ipv6s, *v6.Ipv6Address)
			//}

			for _, pn := range net.PrivateIpAddresses {
				privateIps = append(privateIps, *pn.PrivateIpAddress)
			}

			//tmp.OwnerID = net.OwnerId
		}
		tmp.PrivateIPs = aws.String(strings.Join(privateIps, constants.DefaultDelimiter))
		//tmp.IPv6s = aws.String(strings.Join(ipv6s, constants.DefaultDelimiter))

		var sgNames []string
		var sgIds []string
		for _, sg := range instance.SecurityGroups {
			sgNames = append(sgNames, *sg.GroupName)
			sgIds = append(sgIds, *sg.GroupId)
		}
		tmp.SecurityGroupIDs = aws.String(strings.Join(sgIds, constants.DefaultDelimiter))
		tmp.SecurityGroupNames = aws.String(strings.Join(sgNames, constants.DefaultDelimiter))

		logrus.Tracef("Instance is added: %s", *tmp.InstanceID)
		ch <- tmp
	}

	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			wg.Add(1)
			go f(instance, input)
		}
	}

	wg.Wait()
	close(input)

	result = <-output
	logrus.Debugf("total valid EC2 data count: %d", len(result))

	return result, nil
}

// SetAlias sets alias
func (e *EC2Client) SetAlias(alias *string) {
	e.Alias = alias
}

// GetEC2Instances get all instances in the account
func (e *EC2Client) GetEC2Instances(original []types.Reservation, nextToken *string) ([]types.Reservation, error) {
	result, err := e.Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		NextToken: nextToken,
	})
	if err != nil {
		return nil, err
	}

	original = append(original, result.Reservations...)
	if result.NextToken != nil {
		return e.GetEC2Instances(original, result.NextToken)
	}
	return original, nil
}

// NewEC2Client creates EC2Client resource with ec2 client
func NewEC2Client(cfg aws.Config, helper Helper) (Client, error) {
	return &EC2Client{
		Resource: constants.EC2ResourceName,
		Client:   GetEC2ClientFn(cfg),
		Region:   helper.Region,
	}, nil
}

// GetEC2ClientFn creates ec2 client
func GetEC2ClientFn(cfg aws.Config) *ec2.Client {
	return ec2.NewFromConfig(cfg)
}

// GetAllRegions will returns all regions
func GetAllRegions() ([]string, error) {
	logrus.Debug("Retrieve all regions in the AWS provider")
	svc := ec2.NewFromConfig(GetAwsSession(constants.EmptyString))

	input := &ec2.DescribeRegionsInput{}

	result, err := svc.DescribeRegions(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var regions []string
	for _, region := range result.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}
