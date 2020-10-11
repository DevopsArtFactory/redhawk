package constants

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// DefaultLogLevel is the default global verbosity
	DefaultLogLevel = logrus.WarnLevel

	// DefaultRegion is the default region id
	DefaultRegion = "us-east-1"

	// EmptyString is the empty string
	EmptyString = ""

	// DefaultProvider returns default provider for redhawk
	DefaultProvider = "aws"

	// DefaultOutputFormat is a default format for output
	DefaultOutputFormat = "stdout"

	// DefaultDelimiter is a default delimiter for csv output
	DefaultDelimiter = "|"

	// Resource Name Constants
	// After add resource here, you have to setup `ResourceConfig` in the var section
	EC2ResourceName      = "ec2"
	SGResourceName       = "security_group"
	Route53ResourceName  = "route53"
	S3ResourceName       = "s3"
	RDSResourceName      = "rds"
	IAMResourceName      = "iam"
	IAMUserResourceName  = "iam_user"
	IAMGroupResourceName = "iam_group"
	IAMRoleResourceName  = "iam_role"
)

var (
	// AllAWSRegions is a list of all AWS Region
	AllAWSRegions = []string{
		"eu-north-1",
		"ap-south-1",
		"eu-west-3",
		"eu-west-2",
		"eu-west-1",
		"ap-northeast-2",
		"ap-northeast-1",
		"sa-east-1",
		"ca-central-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"eu-central-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}

	// ValidFormats is a list of valid output format for scan data
	ValidFormats = []string{
		"stdout",
		"csv",
	}
	ResourceConfigs = []ResourceConfig{
		{
			Name:    EC2ResourceName,
			Default: true,
			Global:  false,
		},
		{
			Name:    SGResourceName,
			Default: true,
			Global:  false,
		},
		{
			Name:    Route53ResourceName,
			Default: true,
			Global:  false,
		},
		{
			Name:    S3ResourceName,
			Default: true,
			Global:  false,
		},
		{
			Name:    RDSResourceName,
			Default: true,
			Global:  false,
		},
		{
			Name:    IAMResourceName,
			Default: true,
			Global:  true,
		},
	}
)

// Get Home Directory
func HomeDir() string {
	if h := os.Getenv("HOME"); h != EmptyString {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
