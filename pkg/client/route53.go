package client

import (
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/route53"
	"strings"
)

type Route53Client struct {
	Resource string
	Client   *route53.Route53
}

// GetResourceName returns resource name of client
func (r Route53Client) GetResourceName() string {
	return r.Resource
}

// NewRoute53Client creates EC2Client resource with ec2 client
func NewRoute53Client(helper Helper) (Client, error) {
	session := GetAwsSession()
	return &Route53Client{
		Resource: constants.Route53ResourceName,
		Client:   GetRoute53ClientFn(session, helper.Region, helper.Credentials),
	}, nil
}

// GetRoute53ClientFn creates route53 client
func GetRoute53ClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *route53.Route53 {
	if creds == nil {
		return route53.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return route53.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// Scan scans all data
func (r Route53Client) Scan() ([]resource.Resource, error) {
	var result []resource.Resource

	recordSets, err := r.GetRoute53List()
	if err != nil {
		return nil, err
	}

	for _, rs := range recordSets {
		tmp := resource.Route53Resource{
			ResourceType: aws.String(constants.Route53ResourceName),
		}
		tmp.Name = rs.Name
		tmp.Type = rs.Type

		if rs.AliasTarget != nil {
			tmp.Alias = aws.Bool(true)
			tmp.RouteTo = rs.AliasTarget.DNSName
		}

		if len(rs.ResourceRecords) > 0 {
			var routeTo []string
			for _, rr := range rs.ResourceRecords {
				routeTo = append(routeTo, *rr.Value)
			}
			tmp.RouteTo = aws.String(strings.Join(routeTo, "|"))
			tmp.Alias = aws.Bool(false)
		}

		tmp.TTL = rs.TTL

		result = append(result, tmp)
	}

	return result, nil
}

// GetRoute53List get all record set in the account
func (r Route53Client) GetRoute53List() ([]*route53.ResourceRecordSet, error) {
	hostedZones, err := r.GetRoute53HostedZones()
	if err != nil {
		return nil, err
	}

	var ret []*route53.ResourceRecordSet
	for _, hz := range hostedZones {
		result, err := r.Client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId: hz.Id,
		})

		ret = append(ret, result.ResourceRecordSets...)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

// GetRoute53HostedZones get all hosted zones in the account
func (r Route53Client) GetRoute53HostedZones() ([]*route53.HostedZone, error) {
	result, err := r.Client.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		return nil, err
	}

	return result.HostedZones, nil
}
