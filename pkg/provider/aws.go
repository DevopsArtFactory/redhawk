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
		Provider: "aws",
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
