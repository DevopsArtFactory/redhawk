package resource

import (
	"strings"

	"github.com/jszwec/csvutil"
)

// GetResource returns resource type
func (s S3Resource) GetResource() string {
	return *s.ResourceType
}

// GetHeaders returns headers
func (s S3Resource) GetHeaders() ([]string, error) {
	strSlice, err := s.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	return strings.Split(strSlice[0], ","), nil
}

// TransferToCSV change struct to CSV
func (s S3Resource) TransferToCSV() ([]string, error) {
	strSlice, err := s.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	if len(strSlice) <= 1 {
		return nil, nil
	}

	return strings.Split(strSlice[1], ","), nil
}

// StructToSliceLine returns header and rows
func (s S3Resource) StructToSliceLine() ([]string, error) {
	b, err := csvutil.Marshal([]S3Resource{s})
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(b), "\n")

	return split, nil
}
