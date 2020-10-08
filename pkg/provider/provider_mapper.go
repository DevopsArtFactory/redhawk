package provider

import "github.com/DevopsArtFactory/redhawk/pkg/client"

var (
	providers = map[string]func() Provider {
		"client": NewAWSProvider,
	}
	clientMapperFunc = map[string]func(string, string, string) (client.Client, error) {
		"client": client.CreateResourceClient,
	}
)

