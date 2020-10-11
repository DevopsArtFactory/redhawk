package schema

// Configuration for redhawk
type Config struct {
	// Resource Provider like AWS, GCP etc...
	Provider string `yaml:"provider"`

	// Multi accounts name and role for AWS Provider
	Accounts []Account `yaml:"accounts,omitempty"`

	// List of regions. Default region of provider will be applied if no region specified
	Regions []string `yaml:"regions,omitempty"`

	// List of resources. All resources will be applied if no resources specified
	Resources []Resource `yaml:"resources,omitempty"`
}

// Configuration for assume account for AWS
type Account struct {
	// Custom account identifier
	Name string `yaml:"name,omitempty"`

	// Assume role ARN
	RoleArn string `yaml:"role_arn,omitempty"`
}

// Resource configuration with detailed conditions
type Resource struct {
	// Resource name
	Name   string `yaml:"name"`
	Global bool   `yaml:"global"`
}
