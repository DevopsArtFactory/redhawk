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
