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
