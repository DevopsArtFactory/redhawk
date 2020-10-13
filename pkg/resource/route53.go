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

package resource

import (
	"strings"

	"github.com/jszwec/csvutil"
)

// GetResource returns resource type
func (r Route53Resource) GetResource() string {
	return *r.ResourceType
}

// GetHeaders returns headers
func (r Route53Resource) GetHeaders() ([]string, error) {
	strSlice, err := r.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	return strings.Split(strSlice[0], ","), nil
}

// TransferToCSV change struct to CSV
func (r Route53Resource) TransferToCSV() ([]string, error) {
	strSlice, err := r.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	if len(strSlice) <= 1 {
		return nil, nil
	}

	return strings.Split(strSlice[1], ","), nil
}

// StructToSliceLine returns header and rows
func (r Route53Resource) StructToSliceLine() ([]string, error) {
	b, err := csvutil.Marshal([]Route53Resource{r})
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(b), "\n")

	return split, nil
}
