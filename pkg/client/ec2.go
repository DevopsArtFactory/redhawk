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
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
)

type EC2Client struct {
	Resource string
	Client   *ec2.EC2
}

// GetResourceName returns resource name of client
func (e EC2Client) GetResourceName() string {
	return e.Resource
}

// Scan scans all data
func (e EC2Client) Scan() ([]resource.Resource, error) {
	var wg sync.WaitGroup
	var result []resource.Resource

	logrus.Debugf("Start to scan all ec2 instances")
	reservations, err := e.GetEC2Instances()
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

	f := func(reservation *ec2.Reservation, ch chan resource.EC2Resource) {
		tmp := resource.EC2Resource{
			ResourceType: aws.String(constants.EC2ResourceName),
		}

		logrus.Tracef("Possible valid instances: %d", len(reservation.Instances))
		for _, instance := range reservation.Instances {
			logrus.Tracef("Gathering information about instance: %s", *instance.InstanceId)
			tmp.InstanceID = instance.InstanceId
			tmp.InstanceStatus = instance.State.Name
			tmp.InstanceType = instance.InstanceType
			tmp.AvailabilityZone = instance.Placement.AvailabilityZone
			tmp.LaunchTime = instance.LaunchTime
			tmp.ImageID = instance.ImageId
			tmp.VpcID = instance.VpcId

			if instance.PublicIpAddress != nil {
				tmp.PublicIP = instance.PublicIpAddress
			}

			if instance.KeyName != nil {
				tmp.KeyName = instance.KeyName
			}

			if instance.IamInstanceProfile != nil {
				tmp.IAMInstanceProfile = instance.IamInstanceProfile.Arn
			}

			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					tmp.Name = tag.Value
					break
				}
			}

			var privateIps []string
			var ipv6s []string
			for _, net := range instance.NetworkInterfaces {
				for _, v6 := range net.Ipv6Addresses {
					ipv6s = append(ipv6s, *v6.Ipv6Address)
				}

				for _, pn := range net.PrivateIpAddresses {
					privateIps = append(privateIps, *pn.PrivateIpAddress)
				}

				tmp.OwnerID = net.OwnerId
			}
			tmp.PrivateIPs = aws.String(strings.Join(privateIps, constants.DefaultDelimiter))
			tmp.IPv6s = aws.String(strings.Join(ipv6s, constants.DefaultDelimiter))

			var sgNames []string
			var sgIds []string
			for _, sg := range instance.SecurityGroups {
				sgNames = append(sgNames, *sg.GroupName)
				sgIds = append(sgIds, *sg.GroupId)
			}
			tmp.SecurityGroupIDs = aws.String(strings.Join(sgIds, constants.DefaultDelimiter))
			tmp.SecurityGroupNames = aws.String(strings.Join(sgNames, constants.DefaultDelimiter))
		}

		logrus.Tracef("Instance is added: %s", *tmp.InstanceID)
		ch <- tmp
	}

	logrus.Debugf("Instances found: %d", len(reservations))
	for _, reservation := range reservations {
		wg.Add(1)
		go f(reservation, input)
	}

	wg.Wait()
	close(input)

	result = <-output
	logrus.Debugf("total valid EC2 data count: %d", len(result))

	return result, nil
}

// GetEC2Instances get all instances in the account
func (e EC2Client) GetEC2Instances() ([]*ec2.Reservation, error) {
	result, err := e.Client.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	return result.Reservations, nil
}

// NewEC2Client creates EC2Client resource with ec2 client
func NewEC2Client(helper Helper) (Client, error) {
	session := GetAwsSession()
	return &EC2Client{
		Resource: constants.EC2ResourceName,
		Client:   GetEC2ClientFn(session, helper.Region, helper.Credentials),
	}, nil
}

// GetEC2ClientFn creates ec2 client
func GetEC2ClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *ec2.EC2 {
	if creds == nil {
		return ec2.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return ec2.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// GetAllRegions will returns all regions
func GetAllRegions() ([]string, error) {
	logrus.Debug("Retrieve all regions in the AWS provider")
	svc := ec2.New(GetAwsSession(), &aws.Config{Region: aws.String(constants.DefaultRegion)})

	input := &ec2.DescribeRegionsInput{}

	result, err := svc.DescribeRegions(input)
	if err != nil {
		return nil, err
	}

	var regions []string
	for _, region := range result.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}
