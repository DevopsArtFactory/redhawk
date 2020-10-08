package runner

import (
	"github.com/DevopsArtFactory/redhawk/pkg/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/client"
	"github.com/DevopsArtFactory/redhawk/pkg/color"
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/provider"
	"github.com/DevopsArtFactory/redhawk/pkg/templates"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"text/tabwriter"
)

type Runner struct {
	AWSClient client.Client
	Builder   *builder.Builder
	TotalCount int
}

type Record struct {
	Error error
	Resource string
	Region string
	Data []map[string]interface{}
}

func New() *Runner {
	region := viper.GetString("region")
	if len(region) == 0 {
		region = constants.DefaultRegion
	}

	return &Runner{}
}

// ScanResources retrieves resources in AWS
func (r Runner) ScanResources(out io.Writer) error {
	logrus.Debug("start scanning resources")

	var errors []error
	ch := make(chan Record)
	r.TotalCount = len(r.Builder.Config.Regions) * len(r.Builder.Config.Resources)

	// Region based
	for _, region := range r.Builder.Config.Regions {
		// Create new provider
		prov := provider.CreateProvider(r.Builder.Config.Provider)

		// Resources based
		for _, t := range r.Builder.Config.Resources {
			go func(name, region string) {
				re := Record{
					Error:    nil,
					Resource: t.Name,
					Region:   region,
				}
				c, err := prov.CreateClient(region, name)
				if err != nil {
					re.Error = err
				} else {
					data, err := c.Scan()
					re.Error = err
					re.Data = data
				}

				ch <- re
			}(t.Name, region)
		}
	}

	result := map[string]map[string][]map[string]interface{}{}
	for i := 0; i < r.TotalCount; i++ {
		record := <- ch
		if _, ok := result[record.Resource]; !ok {
			result[record.Resource] = map[string][]map[string]interface{}{}
		}

		if _, ok := result[record.Resource][record.Region]; !ok {
			result[record.Resource][record.Region] = record.Data
		} else {
			result[record.Resource][record.Region] = append(result[record.Resource][record.Region], record.Data...)
		}

		if record.Error != nil {
			errors = append(errors, record.Error)
		}
	}

	if err := PrintScanResult(out, result); err != nil {
		return err
	}

	//for _, err := range errors {
	//	fmt.Println(err.Error())
	//}
	return nil
}

// PrintScanResult prints scan result
func PrintScanResult(out io.Writer, result map[string]map[string][]map[string]interface{}) error {
	funcMap := template.FuncMap{
		"decorate": color.DecorateAttr,
	}

	for resource, ret := range result {
		w := tabwriter.NewWriter(out, 0, 5, 3, ' ', tabwriter.TabIndent)
		t := template.Must(template.New(templates.Templates[resource]["title"]).Funcs(funcMap).Parse(templates.Templates[resource]["template"]))

		err := t.Execute(w, data)
		if err != nil {
			return err
		}
		w.Flush()
	}

	return nil
}