package templates

var (
	Templates = map[string]string{
		"aws": AWSTemplate,
	}
)

// AWSTemplate is a template for aws provider
const AWSTemplate = `PROVIDER: {{ .Provider }}
============================================
{{- if gt (len .Summary.EC2) 0 }}
SERVICE	NAME	ID	STATUS	TYPE	AZ	PUBLIC IP	IPv6	PRIVATE IP	SG NAME	SG ID	KEY	OWNER	IMAGE	LAUNCHED
{{- range $ec2 := .Summary.EC2 }}
EC2	{{ $ec2.Name }}	{{ $ec2.InstanceID }}	{{ $ec2.InstanceStatus }}	{{ $ec2.InstanceType }}	{{ $ec2.AvailabilityZone }}	{{ $ec2.PublicIP }}	{{ $ec2.IPv6s }}	{{ $ec2.PrivateIPs }}	{{ $ec2.SecurityGroupNames }}	{{ $ec2.SecurityGroupIDs }}	{{ $ec2.KeyName }}	{{ $ec2.OwnerID }}	{{ $ec2.ImageID }}	{{ $ec2.LaunchTime }}
{{- end }}
{{- end }}
`
