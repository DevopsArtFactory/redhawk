package schema

import "time"

type AWSResources struct {
	// AWS Resources
	EC2 []EC2Resource `json:"ec2,omitempty"`
}

// EC2 Resource columns
type EC2Resource struct {
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
	IPv6s              []string   `json:"ipv6,omitempty"`
	PrivateIPs         []string   `json:"private_ips,omitempty"`
	SecurityGroupNames []string   `json:"security_group_names,omitempty"`
	SecurityGroupIDs   []string   `json:"security_group_ids,omitempty"`
	LaunchTime         *time.Time `json"launch_time,omitempty"`
}
