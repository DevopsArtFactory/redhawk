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
			tmp[policyIndex] = string(decodedPolicy)
		}

		if resource.GetResource() == constants.Route53ResourceName && len(b[4]) > 0 {
			decodedRouteTo, err := base64.StdEncoding.DecodeString(b[4])
			if err != nil {
				return nil, err
			}
			b[4] = string(decodedRouteTo)
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
