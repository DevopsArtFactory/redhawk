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

package printer

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
)

type CSVPrinter struct {
	Out      *os.File
	Provider string
	Data     map[string][][]string
}

func NewCSVPrinter() Printer {
	return CSVPrinter{}
}

// SetData sets data
func (c CSVPrinter) SetData(provider string, d []resource.Resource) (Printer, error) {
	ret := map[string][][]string{}
	for _, resource := range d {
		rt := resource.GetResource()
		_, ok := ret[rt]
		if !ok {
			header, err := resource.GetHeaders()
			if err != nil {
				return nil, err
			}
			ret[rt] = [][]string{
				header,
			}
		}

		b, err := resource.TransferToCSV()
		if err != nil {
			return nil, err
		}

		tmp := b
		policyIndex := 6
		if resource.GetResource() == constants.S3ResourceName && len(b[policyIndex]) > 0 {
			decodedPolicy, err := base64.StdEncoding.DecodeString(b[policyIndex])
			if err != nil {
				return nil, err
			}
			tmp = b[:policyIndex]
			tmp = append(tmp, string(decodedPolicy))
		}

		routeToIndex := 4
		if resource.GetResource() == constants.Route53ResourceName && len(b[routeToIndex]) > 0 {
			decodedRouteTo, err := base64.StdEncoding.DecodeString(b[routeToIndex])
			if err != nil {
				return nil, err
			}
			b[routeToIndex] = string(decodedRouteTo)
			tmp = b
		}

		ret[rt] = append(ret[rt], tmp)
	}

	c.Data = ret
	c.Provider = provider

	return c, nil
}

// Print shows data to Standard Out
func (c CSVPrinter) Print() error {
	now := time.Now().Unix()
	for key, dataList := range c.Data {
		filePath := getRandomFilePath(c.Provider, key, now)
		if !tools.FileExists(filePath) {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("cannot create file: %s", filePath)
			}
			w := csv.NewWriter(f)
			defer f.Close()

			for _, data := range dataList {
				if err := w.Write(data); err != nil {
					return err
				}
			}

			w.Flush()

			if err := w.Error(); err != nil {
				return err
			}
		}
	}

	return nil
}

// getRandomFilePath creates filename for csv
func getRandomFilePath(provider, key string, now int64) string {
	return fmt.Sprintf("%s-%d-%s.csv", provider, now, key)
}
