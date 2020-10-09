package provider

import (
	"fmt"

	"github.com/DevopsArtFactory/redhawk/pkg/client"
)

type Provider interface {
	CreateClient(string, string) (client.Client, error)
	GetProvider() string
}

// CreateProvider creates new provider for redhawk
func CreateProvider(provider string) (Provider, error) {
	return ChooseProvider(provider)
}

// ChooseProvider select provider with specified provider key
func ChooseProvider(provider string) (Provider, error) {
	f, ok := providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider is not supported: %s", provider)
	}
	return f(), nil
}
