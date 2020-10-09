package provider

var (
	providers = map[string]func() Provider{
		"aws": NewAWSProvider,
	}
)
