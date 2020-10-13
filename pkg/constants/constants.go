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

package constants

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// DefaultLogLevel is the default global verbosity
	DefaultLogLevel = logrus.InfoLevel

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

	// DefaultRegionVariable is the default region id
	DefaultRegionVariable = "AWS_DEFAULT_REGION"

	// StringText is "string"
	StringText = "string"

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

	ResourceGlobal = map[string]bool{
		EC2ResourceName:     false,
		SGResourceName:      false,
		Route53ResourceName: false,
		S3ResourceName:      false,
		RDSResourceName:     false,
		IAMResourceName:     true,
	}

	// AWSCredentialsPath is the file path of aws credentials
	AWSCredentialsPath = HomeDir() + "/.aws/credentials"

	// AWSConfigPath is the file path of aws config
	AWSConfigPath = HomeDir() + "/.aws/config"

	ResourceConfigs = []ResourceConfig{
		{
			Name:    EC2ResourceName,
			Default: true,
		},
		{
			Name:    SGResourceName,
			Default: true,
		},
		{
			Name:    Route53ResourceName,
			Default: true,
		},
		{
			Name:    S3ResourceName,
			Default: true,
		},
		{
			Name:    RDSResourceName,
			Default: true,
		},
		{
			Name:    IAMResourceName,
			Default: true,
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
