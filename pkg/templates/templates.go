package templates

var (
	Templates = map[string]string{
		"aws": AWSTemplate,
	}
)

// AWSTemplate is a template for aws provider
const AWSTemplate = `PROVIDER: {{ .Provider }}
{{- range $key, $val := .Summary }}
  {{- if eq $key "ec2" }}
    {{- if gt (len $val) 0 }}
	  {{- if $.Detail }}
==============================================
SERVICE	NAME	ID	STATUS	TYPE	AZ	PUBLIC IP	IPv6	PRIVATE IP	SG NAME	SG ID	KEY	OWNER	IMAGE	LAUNCHED
	    {{- range $ec2 := $val }}
EC2	{{ $ec2.Name }}	{{ $ec2.InstanceID }}	{{ $ec2.InstanceStatus }}	{{ $ec2.InstanceType }}	{{ $ec2.AvailabilityZone }}	{{ $ec2.PublicIP }}	{{ $ec2.IPv6s }}	{{ $ec2.PrivateIPs }}	{{ $ec2.SecurityGroupNames }}	{{ $ec2.SecurityGroupIDs }}	{{ $ec2.KeyName }}	{{ $ec2.OwnerID }}	{{ $ec2.ImageID }}	{{ $ec2.LaunchTime }}
	    {{- end }}
	  {{- else }}
==============================================
SERVICE	NAME	ID	STATUS	TYPE	AZ	LAUNCHED
	    {{- range $ec2 := $val }}
EC2	{{ $ec2.Name }}	{{ $ec2.InstanceID }}	{{ $ec2.InstanceStatus }}	{{ $ec2.InstanceType }}	{{ $ec2.AvailabilityZone }}	{{ $ec2.LaunchTime }}
	    {{- end }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "security_group" }}
    {{- if gt (len $val) 0 }}

==============================================
SERVICE	NAME	ID	VPC	OWNER	INBOUND	OUTBOUND	DESCRIPTION
	  {{- range $sg := $val }}
SG	{{ $sg.Name }}	{{ $sg.ID }}	{{ $sg.VpcID }}	{{ $sg.Owner }}	{{ $sg.InboundCount }}	{{ $sg.OutboundCount }}	{{ $sg.Description }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "route53" }}
    {{- if gt (len $val) 0 }}

==============================================
SERVICE	NAME	TYPE	ALIAS	TARGET	TTL
	  {{- range $route53 := $val }}
Route53	{{ $route53.Name }}	{{ $route53.Type }}	{{ $route53.Alias }}	{{ $route53.RouteTo }}	{{ $route53.TTL }}
	  {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
`
