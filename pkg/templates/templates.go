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
SERVICE	STATUS	NAME	ID	TYPE	AZ	Region	SG_NAME	SG_ID	SUBNET_ID	PUBLIC_IP	PRIVATE_IP	IMAGE	VPC_ID	KEY	LAUNCHED
	    {{- range $ec2 := $val }}
EC2	{{ format $ec2.InstanceStatus }}	{{ format $ec2.Name }}	{{ format $ec2.InstanceID }}	{{ format $ec2.InstanceType }}	{{ format $ec2.AvailabilityZone }}	{{ format $ec2.RegionName }}	{{ format $ec2.SecurityGroupNames }}	{{ format $ec2.SecurityGroupIDs }}	{{ format $ec2.SubnetID }}	{{ format $ec2.PublicIP }}	{{ format $ec2.PrivateIPs }}	{{ format $ec2.ImageID }}	{{ format $ec2.VpcID }}	{{ format $ec2.KeyName }}	{{ format $ec2.LaunchTime }}
	    {{- end }}
	  {{- else }}
==============================================
SERVICE	STATUS	NAME	ID	TYPE	AZ	LAUNCHED
	    {{- range $ec2 := $val }}
EC2	{{ format $ec2.InstanceStatus }}	{{ format $ec2.Name }}	{{ format $ec2.InstanceID }}	{{ format $ec2.InstanceType }}	{{ format $ec2.AvailabilityZone }}	{{ format $ec2.LaunchTime }}
	    {{- end }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "security_group" }}
    {{- if gt (len $val) 0 }}

==============================================
SERVICE	NAME	ID	VPC	OWNER	INBOUND	OUTBOUND	DESCRIPTION
	  {{- range $sg := $val }}
SG	{{ $sg.Name }}	{{ format $sg.ID }}	{{ format $sg.VpcID }}	{{ format $sg.Owner }}	{{ format $sg.InboundCount }}	{{ format $sg.OutboundCount }}	{{ format $sg.Description }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "route53" }}
    {{- if gt (len $val) 0 }}

==============================================
SERVICE	NAME	TYPE	ALIAS	TARGET	TTL
	  {{- range $route53 := $val }}
Route53	{{ $route53.Name }}	{{ format $route53.Type }}	{{ format $route53.Alias }}	{{ format $route53.RouteTo }}	{{ format $route53.TTL }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "s3" }}
    {{- if gt (len $val) 0 }}

==============================================
SERVICE	NAME	REGION	LOGGING_ENABLED	LOGGING_BUCKET	CREATED
	  {{- range $s3 := $val }}
S3	{{ $s3.Bucket }}	{{ format $s3.Region }}	{{ format $s3.LoggingEnabled }}	{{ format $s3.LoggingBucket }}	{{ format $s3.Created }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "rds" }}
    {{- if gt (len $val) 0 }}
	  {{- if $.Detail }}
==============================================
SERVICE	IDENTIFIER	ROLE	ENGINE	VERSION	SIZE	STATUS	AZ	STORAGE	OPTION_GROUPS	PARAMETER_GROUPS	SUBNET_GROUP	CREATED
	    {{- range $rds := $val }}
RDS	{{ $rds.RDSIdentifier }}	{{ format $rds.Role }}	{{ format $rds.Engine }}	{{ format $rds.EngineVersion }}	{{ format $rds.Size }}	{{ format $rds.Status }}	{{ format $rds.AvailabilityZone }}	{{ format $rds.StorageType }}	{{ format $rds.OptionGroup }}	{{ format $rds.ParameterGroup }}	{{ format $rds.DBSubnet }}	{{ format $rds.Created }}
	    {{- end }}
	  {{- else }}
==============================================
SERVICE	IDENTIFIER	ROLE	ENGINE	SIZE	STATUS	AZ	CREATED
	    {{- range $rds := $val }}
RDS	{{ $rds.RDSIdentifier }}	{{ format $rds.Role }}	{{ format $rds.Engine }}	{{ format $rds.Size }}	{{ format $rds.Status }}	{{ format $rds.AvailabilityZone }}	{{ format $rds.Created }}
	    {{- end }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "iam_user" }}
    {{- if gt (len $val) 0 }}
==============================================
SERVICE	NAME	MFA	GROUP_COUNT	ACCESS_KEY_LAST_USED	CREATED
	  {{- range $iamUser := $val }}
IAM_USER	{{ $iamUser.UserName }}	{{ format $iamUser.MFA }}	{{ format $iamUser.GroupCount }}	{{ format $iamUser.AccessKeyLastUsed }}	{{ format $iamUser.UserCreated }}
	  {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "iam_group" }}
    {{- if gt (len $val) 0 }}
==============================================
SERVICE	NAME	USER_COUNT	USERS	GROUP_POLICIES
	  {{- range $iamGroup := $val }}
IAM_GROUP	{{ $iamGroup.GroupName }}	{{ format $iamGroup.UserCount }}	{{ format $iamGroup.Users }}	{{ format $iamGroup.GroupPolicies }}
      {{- end }}
    {{- end }}
  {{- end }}

  {{- if eq $key "iam_role" }}
    {{- if gt (len $val) 0 }}
==============================================
SERVICE	NAME	TRUST_ENTITIES	ROLE_LAST_ACTIVITY
	  {{- range $iamRole := $val }}
IAM_ROLE	{{ format $iamRole.RoleName }}	{{ format $iamRole.TrustedEntities }}	{{ format $iamRole.RoleLastActivity }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
`

const HelperTemplates = `redhawk list command

You are currently use {{ .Account }}. If you run redhawk, then you will get the resources based on this account.

Usage:
  redhawk list [flags] [options]

Example:
  # Find resources in detailed version
  - redhawk list --resources=ec2 --detail

  # Find ec2 resource only and region set to Seoul
  - redhawk list --resources=ec2 --region=ap-northeast-2

  # Find ec2,iam resources and region set to all(all region)
  - redhawk list --resources=ec2,iam --all

  # Find ec2,iam resources and region set to us-west-2 and output set to csv
  - redhawk list --resources=ec2,iam --region=us-west-2 -o csv

Options:
      --resources='': [Required]Resource list of provider for dynamic search(Delimiter: comma)
  -A, --all=false: Apply all regions of provider for command
      --config='': configuration file path for scanning resources
      --detail=false: detailed options for scanning
  -o, --output='stdout': detailed options for scanning
  -r, --region='': Run command to specific region
`
