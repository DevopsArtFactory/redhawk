package constants

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// DefaultLogLevel is the default global verbosity
	DefaultLogLevel = logrus.WarnLevel

	// DefaultRegion is the default region id
	DefaultRegion = "ap-northeast-2"

	// EmptyString is the empty string
	EmptyString = ""

	// DefaultProvider returns default provider for redhawk
	DefaultProvider = "client"

	// DefaultProfile is the default client profile
	DefaultProfile = "default"

	// InfoLogLevel is the info level verbosity
	InfoLogLevel = logrus.InfoLevel
)

var (
	// AllRegions is a list of all AWS Region
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
)

// Get Home Directory
func HomeDir() string {
	if h := os.Getenv("HOME"); h != EmptyString {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
