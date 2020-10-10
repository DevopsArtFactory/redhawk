package resource

import (
	"strings"

	"github.com/jszwec/csvutil"
)

// GetResource returns resource type
func (e EC2Resource) GetResource() string {
	return *e.ResourceType
}

// GetHeaders returns headers
func (e EC2Resource) GetHeaders() ([]string, error) {
	strSlice, err := e.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	return strings.Split(strSlice[0], ","), nil
}

// TransferToCSV change struct to CSV
func (e EC2Resource) TransferToCSV() ([]string, error) {
	strSlice, err := e.StructToSliceLine()
	if err != nil {
		return nil, err
	}

	if len(strSlice) <= 1 {
		return nil, nil
	}

	return strings.Split(strSlice[1], ","), nil
}

// StructToSliceLine returns header and rows
func (e EC2Resource) StructToSliceLine() ([]string, error) {
	b, err := csvutil.Marshal([]EC2Resource{e})
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(b), "\n")

	return split, nil
}
