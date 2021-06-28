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
)

var RegionNameMapping = map[string]string{
	"eu-north-1":     "Stockholm",
	"ap-south-1":     "Mumbai",
	"eu-west-3":      "Paris",
	"eu-west-2":      "London",
	"eu-west-1":      "Ireland",
	"ap-northeast-2": "Seoul",
	"ap-northeast-1": "Tokyo",
	"sa-east-1":      "São Paulo",
	"ca-central-1":   "Canada(Central)",
	"ap-southeast-1": "Singapore",
	"ap-southeast-2": "Sydney",
	"eu-central-1":   "Frankfurt",
	"us-east-1":      "N. Virginia",
	"us-east-2":      "Ohio",
	"us-west-1":      "N. California",
	"us-west-2":      "Oregon",
}

// GetRegionName returns exact name of region
func GetRegionName(regionID string) (*string, error) {
	var name string
	var ok bool

	if name, ok = RegionNameMapping[regionID]; !ok {
		return nil, fmt.Errorf("region is not valid: %s", regionID)
	}

	return &name, nil
}
