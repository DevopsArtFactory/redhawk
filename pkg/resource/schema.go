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

package resource

import "time"

type Resource interface {
	GetHeaders() ([]string, error)
	TransferToCSV() ([]string, error)
	GetResource() string
	StructToSliceLine() ([]string, error)
}

type Resources struct {
	// Provider
	Provider string

	// Resources
	Resources []Resource
}

// EC2 Resource columns
type EC2Resource struct {
	ResourceType       *string    `json:"resource_type,omitempty"`
	InstanceStatus     *string    `json:"instance_status,omitempty"`
	AccountAlias       *string    `json:"account_alias,omitempty"`
	Name               *string    `json:"name,omitempty"`
	InstanceID         *string    `json:"instance_id,omitempty"`
	InstanceType       *string    `json:"instance_type,omitempty"`
	AvailabilityZone   *string    `json:"availability_zone,omitempty"`
	RegionName         *string    `json:"region_name,omitempty"`
	SecurityGroupNames *string    `json:"security_group_names,omitempty"`
	SecurityGroupIDs   *string    `json:"security_group_ids,omitempty"`
	SubnetID           *string    `json:"subnet_id,omitempty"`
	PublicIP           *string    `json:"public_ip,omitempty"`
	PrivateIPs         *string    `json:"private_ips,omitempty"`
	ImageID            *string    `json:"image_id,omitempty"`
	VpcID              *string    `json:"vpc_id,omitempty"`
	KeyName            *string    `json:"key_name,omitempty"`
	LaunchTime         *time.Time `json:"launch_time,omitempty"`

	//OwnerID            *string    `json:"owner_id,omitempty"`
	//IPv6s              *string    `json:"ipv6,omitempty"`
	//IAMInstanceProfile *string    `json:"iam_instance_profile,omitempty"`
}

// Security Group Resource columns
type SGResource struct {
	ResourceType  *string `json:"resource_type,omitempty"`
	Name          *string `json:"name,omitempty"`
	ID            *string `json:"id,omitempty"`
	VpcID         *string `json:"vpc_id,omitempty"`
	Owner         *string `json:"owner,omitempty"`
	InboundCount  *int    `json:"inbound_count,omitempty"`
	OutboundCount *int    `json:"outbound_count,omitempty"`
	Description   *string `json:"description,omitempty"`
}

// Route53 Resource columns
type Route53Resource struct {
	ResourceType *string `json:"resource_type,omitempty"`
	Name         *string `json:"name,omitempty"`
	Type         *string `json:"type,omitempty"`
	Alias        *bool   `json:"alias,omitempty"`
	RouteTo      *string `json:"route_to,omitempty"`
	TTL          *int64  `json:"ttl,omitempty"`
}

type S3Resource struct {
	ResourceType   *string    `json:"resource_type,omitempty"`
	Bucket         *string    `json:"bucket,omitempty"`
	Region         *string    `json:"region,omitempty"`
	LoggingEnabled *bool      `json:"logging_enabled,omitempty"`
	LoggingBucket  *string    `json:"logging_bucket,omitempty"`
	Created        *time.Time `json:"created,omitempty"`
	Policy         *string    `json:"policy,omitempty"`
}

type RDSResource struct {
	ResourceType     *string    `json:"resource_type,omitempty"`
	RDSIdentifier    *string    `json:"rds_identifier,omitempty"`
	Role             *string    `json:"role,omitempty"`
	Engine           *string    `json:"engine,omitempty"`
	EngineVersion    *string    `json:"engine_version,omitempty"`
	Region           *string    `json:"region,omitempty"`
	AvailabilityZone *string    `json:"availability_zone,omitempty"`
	Size             *string    `json:"size,omitempty"`
	Status           *string    `json:"status,omitempty"`
	VPC              *string    `json:"vpc,omitempty"`
	StorageType      *string    `json:"storage_type,omitempty"`
	SecurityGroup    *string    `json:"security_group,omitempty"`
	DBSubnet         *string    `json:"db_subnet,omitempty"`
	ParameterGroup   *string    `json:"parameter_group,omitempty"`
	OptionGroup      *string    `json:"option_group,omitempty"`
	Created          *time.Time `json:"created,omitempty"`
}

type IAMUserResource struct {
	ResourceType      *string    `json:"resource_type,omitempty"`
	UserName          *string    `json:"user_name,omitempty"`
	AccessKeyAge      *string    `json:"access_key_age,omitempty"`
	PasswordAge       *string    `json:"password_age,omitempty"`
	UserLastActivity  *string    `json:"user_last_activity,omitempty"`
	MFA               *string    `json:"mfa,omitempty"`
	GroupCount        *int       `json:"group_count,omitempty"`
	ConsoleLastLogin  *time.Time `json:"console_last_login,omitempty"`
	AccessKeyLastUsed *time.Time `json:"access_key_last_usec,omitempty"`
	UserCreated       *time.Time `json:"created,omitempty"`
}

type IAMGroupResource struct {
	ResourceType  *string `json:"resource_type,omitempty"`
	GroupName     *string `json:"group_name,omitempty"`
	Users         *string `json:"users,omitempty"`
	UserCount     *int    `json:"user_count,omitempty"`
	GroupPolicies *string `json:"group_policies,omitempty"`
}

type IAMRoleResource struct {
	ResourceType     *string    `json:"resource_type,omitempty"`
	RoleName         *string    `json:"role_name,omitempty"`
	TrustedEntities  *string    `json:"trusted_entities,omitempty"`
	RoleLastActivity *time.Time `json:"role_last_activity,omitempty"`
}
