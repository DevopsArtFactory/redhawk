package client

import (
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SGClient struct {
	Resource string
	Client   *ec2.EC2
}

// GetResourceName returns resource name of client
func (s SGClient) GetResourceName() string {
	return s.Resource
}

// NewSGClient creates EC2Client resource with ec2 client
func NewSGClient(helper Helper) (Client, error) {
	session := GetAwsSession()
	return &SGClient{
		Resource: constants.SGResourceName,
		Client:   GetEC2ClientFn(session, helper.Region, helper.Credentials),
	}, nil
}

// Scan scans all data
func (s SGClient) Scan() ([]resource.Resource, error) {
	var result []resource.Resource

	securityGroups, err := s.GetSGList()
	if err != nil {
		return nil, err
	}
	for _, sg := range securityGroups {
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

		result = append(result, tmp)
	}

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
