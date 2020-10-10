package client

import (
	"fmt"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Client interface {
	GetResourceName() string
	Scan() ([]resource.Resource, error)
}

type Helper struct {
	Provider string
	Resource string
	Region   string

	// AWS Helper config
	Credentials *credentials.Credentials
}

// ChooseResourceClient selects resource client from the list
func ChooseResourceClient(resource string, h Helper) (Client, error) {
	f, ok := clientMapper[resource]
	if !ok {
		return nil, fmt.Errorf("client does not support: %s", resource)
	}

	c, err := f(h)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// CreateResourceClient creates a new client fro redhawk
func CreateResourceClient(provider, resource, region string) (Client, error) {
	h := Helper{
		Provider: provider,
		Region:   region,
	}

	return ChooseResourceClient(resource, h)
}
