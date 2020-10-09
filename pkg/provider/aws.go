package provider

import (
	"github.com/DevopsArtFactory/redhawk/pkg/client"
)

// AWS provider Struct
type AWSProvider struct {
	Provider string
}

// NewAWSProvider creates AWS Provider
func NewAWSProvider() Provider {
	return AWSProvider{
		Provider: "client",
	}
}

// CreateClient creates a new resource-specific client
func (a AWSProvider) CreateClient(region string, resource string) (client.Client, error) {
	return client.CreateResourceClient(a.Provider, resource, region)
}

// GetProvider returns provider
func (a AWSProvider) GetProvider() string {
	return a.Provider
}
