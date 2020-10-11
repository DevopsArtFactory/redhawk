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
	Name               *string    `json:"name,omitempty"`
	InstanceID         *string    `json:"instance_id,omitempty"`
	InstanceStatus     *string    `json:"instance_status,omitempty"`
	InstanceType       *string    `json:"instance_type,omitempty"`
	AvailabilityZone   *string    `json:"availability_zone,omitempty"`
	ImageID            *string    `json:"image_id,omitempty"`
	PublicIP           *string    `json:"public_ip,omitempty"`
	KeyName            *string    `json:"key_name,omitempty"`
	IAMInstanceProfile *string    `json:"iam_instance_profile,omitempty"`
	VpcID              *string    `json:"vpc_id,omitempty"`
	OwnerID            *string    `json:"owner_id,omitempty"`
	IPv6s              *string    `json:"ipv6,omitempty"`
	PrivateIPs         *string    `json:"private_ips,omitempty"`
	SecurityGroupNames *string    `json:"security_group_names,omitempty"`
	SecurityGroupIDs   *string    `json:"security_group_ids,omitempty"`
	LaunchTime         *time.Time `json:"launch_time,omitempty"`
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
