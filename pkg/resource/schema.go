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
	OwnerID            *string    `json"owner_id,omitempty"`
	IPv6s              *string    `json:"ipv6,omitempty"`
	PrivateIPs         *string    `json:"private_ips,omitempty"`
	SecurityGroupNames *string    `json:"security_group_names,omitempty"`
	SecurityGroupIDs   *string    `json:"security_group_ids,omitempty"`
	LaunchTime         *time.Time `json"launch_time,omitempty"`
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
