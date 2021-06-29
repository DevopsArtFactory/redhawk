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

package client

import (
	"fmt"

	"github.com/DevopsArtFactory/redhawk/pkg/resource"
)

type Client interface {
	GetResourceName() string
	SetAlias(*string)
	Scan() ([]resource.Resource, error)
}

type Helper struct {
	Provider string
	Resource string
	Region   string
}

// ChooseResourceClient selects resource client from the list
func ChooseResourceClient(resource string, h Helper) (Client, error) {
	var err error
	cfg := GetAwsSession(h.Region)

	// Get Account alias
	alias, err := GetAccountAlias(cfg)
	if err != nil {
		return nil, err
	}

	f, ok := clientMapper[resource]
	if !ok {
		return nil, fmt.Errorf("client does not support: %s", resource)
	}

	c, err := f(cfg, h)
	if err != nil {
		return nil, err
	}

	c.SetAlias(&alias)

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
