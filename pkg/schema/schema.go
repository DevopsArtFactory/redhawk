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
	Name string `yaml:"name"`

	// Whether or not it is a global resource or not
	Global bool `yaml:"global"`
}
