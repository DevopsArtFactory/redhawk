package client

import (
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
	"strings"
)

type RDSClient struct {
	Resource string
	Client   *rds.RDS
}

// GetResourceName returns resource name of client
func (r RDSClient) GetResourceName() string {
	return r.Resource
}

// NewRDSClient creates a RDSClient
func NewRDSClient(helper Helper) (Client, error) {
	session := GetAwsSession()
	return &RDSClient{
		Resource: constants.RDSResourceName,
		Client:   GetRDSClientFn(session, helper.Region, helper.Credentials),
	}, nil
}

// GetRDSClientFn creates rds client
func GetRDSClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *rds.RDS {
	if creds == nil {
		return rds.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return rds.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// Scan scans all data
func (r RDSClient) Scan() ([]resource.Resource, error) {
	var result []resource.Resource

	clusters, err := r.GetRDSClusterList()
	if err != nil {
		return nil, err
	}

	logrus.Debugf("RDS clusters found: %d", len(clusters))
	for _, cluster := range clusters {
		for _, dbMember := range cluster.DBClusterMembers {
			tmp := resource.RDSResource{
				ResourceType: aws.String(constants.RDSResourceName),
			}

			tmp.RDSIdentifier = dbMember.DBInstanceIdentifier
			role := "reader"
			if *dbMember.IsClusterWriter {
				role = "writer"
			}

			tmp.Role = aws.String(role)
			tmp.Engine = cluster.Engine
			tmp.EngineVersion = cluster.EngineVersion

			dbInfo, err := r.GetRDSInfo(*dbMember.DBInstanceIdentifier)
			if err != nil {
				return nil, err
			}

			tmp.AvailabilityZone = dbInfo.AvailabilityZone
			tmp.Size = dbInfo.DBInstanceClass
			tmp.Status = dbInfo.DBInstanceStatus
			tmp.VPC = dbInfo.DBSubnetGroup.VpcId
			tmp.StorageType = dbInfo.StorageType
			tmp.DBSubnet = dbInfo.DBSubnetGroup.DBSubnetGroupName
			tmp.Created = dbInfo.InstanceCreateTime

			var sgList []string
			for _, vpcSgID := range dbInfo.VpcSecurityGroups {
				sgList = append(sgList, *vpcSgID.VpcSecurityGroupId)
			}
			tmp.SecurityGroup = aws.String(strings.Join(sgList, "|"))

			var parameterGroups []string
			for _, pg := range dbInfo.DBParameterGroups {
				parameterGroups = append(parameterGroups, *pg.DBParameterGroupName)
			}
			tmp.ParameterGroup = aws.String(strings.Join(parameterGroups, "|"))

			var optionGroups []string
			for _, og := range dbInfo.OptionGroupMemberships {
				optionGroups = append(optionGroups, *og.OptionGroupName)
			}
			tmp.OptionGroup = aws.String(strings.Join(optionGroups, "|"))

			logrus.Tracef("Add new rds instance: %s / %s", *tmp.RDSIdentifier, *tmp.Role)
			result = append(result, tmp)
		}

	}

	return result, nil
}

// GetRDSClusterList returns all DB clusters list in the account
func (r RDSClient) GetRDSClusterList() ([]*rds.DBCluster, error) {
	result, err := r.Client.DescribeDBClusters(&rds.DescribeDBClustersInput{})
	if err != nil {
		return nil, err
	}

	return result.DBClusters, nil
}

// GetRDSInfo returns DB instance information
func (r RDSClient) GetRDSInfo(identifier string) (*rds.DBInstance, error) {
	result, err := r.Client.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(identifier),
	})
	if err != nil {
		return nil, err
	}

	return result.DBInstances[0], nil
}
