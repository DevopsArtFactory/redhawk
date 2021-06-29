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
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// GetAwsSession creates new session for AWS
func GetAwsSession(region string) aws.Config {
	var optFunc config.LoadOptionsFunc
	if len(region) > 0 {
		optFunc = config.WithRegion(region)
	}
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		optFunc,
	)
	return cfg
}
