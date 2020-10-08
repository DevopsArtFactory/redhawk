package templates

var (
	Templates = map[string]map[string]string {
		"ec2": {
			"title": "EC2 instance status",
			"template": EC2Template,
		},
	}
)

// EC2Template
const EC2Template = `{{decorate "ec2" ""}}{{decorate "underline bold" "EC2"}}
{{- if eq (len .Summary) 0 }}
 No rule exists
{{- else }}
ID	TYPE	ACTION	PRIORITY	DATESET COUNT
{{- range $rule := .Summary.Rules }}
{{ $rule.RuleID }}	{{ $rule.Type }}	{{ $rule.ActionType }}	{{ $rule.Priority }}	{{ len $rule.IPDataSet }}
{{- end }}
{{- end }}

{{decorate "IP Set" ""}}{{decorate "underline bold" "IP Set Details"}}
{{- range $rule := .Summary.Rules }}
{{- if gt (len $rule.IPDataSet) 0 }}
{{- range $ipset := $rule.IPDataSet }}
{{- if eq (len $ipset.IPList) 0 }}
No IP is registered
{{- else }}
ID	Count
{{ $ipset.ID }}	{{ len $ipset.IPList }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
`

const IPSearchResultTemplate = `{{decorate "result" ""}}{{decorate "underline bold" "Result"}}
{{- if eq (len .Summary) 0 }}
 No result exists
{{- else }}
IP	IPSetID	Result
{{- range $result := .Summary }}
{{ $result.IP }}	{{ $result.IPSetID }}	{{ $result.Result }}
{{- end }}
{{- end }}
`
