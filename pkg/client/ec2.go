package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/schema"
)

type EC2Client struct {
	Resource string
	Client   *ec2.EC2
}

func (e EC2Client) GetResourceName() string {
	return e.Resource
}

// Scan scans all data
func (e EC2Client) Scan() (*schema.AWSResources, error) {
	var result []schema.EC2Resource

	reservations, err := e.GetEC2Instances()
	if err != nil {
		return nil, err
	}
	for _, reservation := range reservations {
		tmp := schema.EC2Resource{}
		for _, instance := range reservation.Instances {
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

			privateIps := []string{}
			ipv6s := []string{}
			for _, net := range instance.NetworkInterfaces {
				for _, v6 := range net.Ipv6Addresses {
					ipv6s = append(ipv6s, *v6.Ipv6Address)
				}

				for _, pn := range net.PrivateIpAddresses {
					privateIps = append(privateIps, *pn.PrivateIpAddress)
				}

				tmp.OwnerID = net.OwnerId
			}
			tmp.PrivateIPs = privateIps
			tmp.IPv6s = ipv6s

			sgNames := []string{}
			sgIds := []string{}
			for _, sg := range instance.SecurityGroups {
				sgNames = append(sgNames, *sg.GroupName)
				sgIds = append(sgIds, *sg.GroupId)
			}
			tmp.SecurityGroupIDs = sgIds
			tmp.SecurityGroupNames = sgNames
		}

		result = append(result, tmp)
	}

	scanResult := schema.AWSResources{
		EC2: result,
	}

	return &scanResult, nil
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
		Resource: "ec2",
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
