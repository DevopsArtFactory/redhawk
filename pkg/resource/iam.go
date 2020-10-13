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

/*
	IAM USER
*/
// GetResource returns resource type
func (i IAMUserResource) GetResource() string {
	return *i.ResourceType
}

// GetHeaders returns headers
func (i IAMUserResource) GetHeaders() ([]string, error) {
	strSlice, err := i.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	return strings.Split(strSlice[0], ","), nil
}

// TransferToCSV change struct to CSV
func (i IAMUserResource) TransferToCSV() ([]string, error) {
	strSlice, err := i.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	if len(strSlice) <= 1 {
		return nil, nil
	}

	return strings.Split(strSlice[1], ","), nil
}

// StructToSliceLine returns header and rows
func (i IAMUserResource) StructToSliceLine() ([]string, error) {
	b, err := csvutil.Marshal([]IAMUserResource{i})
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(b), "\n")

	return split, nil
}

/*
	IAM GROUP
*/
// GetResource returns resource type
func (i IAMGroupResource) GetResource() string {
	return *i.ResourceType
}

// GetHeaders returns headers
func (i IAMGroupResource) GetHeaders() ([]string, error) {
	strSlice, err := i.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	return strings.Split(strSlice[0], ","), nil
}

// TransferToCSV change struct to CSV
func (i IAMGroupResource) TransferToCSV() ([]string, error) {
	strSlice, err := i.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	if len(strSlice) <= 1 {
		return nil, nil
	}

	return strings.Split(strSlice[1], ","), nil
}

// StructToSliceLine returns header and rows
func (i IAMGroupResource) StructToSliceLine() ([]string, error) {
	b, err := csvutil.Marshal([]IAMGroupResource{i})
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(b), "\n")

	return split, nil
}

/*
	IAM ROLE
*/
// GetResource returns resource type
func (i IAMRoleResource) GetResource() string {
	return *i.ResourceType
}

// GetHeaders returns headers
func (i IAMRoleResource) GetHeaders() ([]string, error) {
	strSlice, err := i.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	return strings.Split(strSlice[0], ","), nil
}

// TransferToCSV change struct to CSV
func (i IAMRoleResource) TransferToCSV() ([]string, error) {
	strSlice, err := i.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	if len(strSlice) <= 1 {
		return nil, nil
	}

	return strings.Split(strSlice[1], ","), nil
}

// StructToSliceLine returns header and rows
func (i IAMRoleResource) StructToSliceLine() ([]string, error) {
	b, err := csvutil.Marshal([]IAMRoleResource{i})
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(b), "\n")

	return split, nil
}
