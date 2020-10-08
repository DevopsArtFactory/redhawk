package provider

import "github.com/DevopsArtFactory/redhawk/pkg/client"

type Provider interface {
	CreateClient(string, string) (client.Client, error)
	GetProvider() string
}

// CreateProvider creates new provider for redhawk
func CreateProvider(provider string) Provider {
	return ChooseProvider(provider)
}

// ChooseProvider select provider with specified provider key
func ChooseProvider(provider string) Provider {
	return providers[provider]()
}