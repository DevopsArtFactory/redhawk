package runner

import (
	"html/template"
	"io"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/redhawk/pkg/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/client"
	"github.com/DevopsArtFactory/redhawk/pkg/color"
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/provider"
	"github.com/DevopsArtFactory/redhawk/pkg/schema"
	"github.com/DevopsArtFactory/redhawk/pkg/templates"
)

type Runner struct {
	AWSClient  client.Client
	Builder    *builder.Builder
	TotalCount int
}

type Record struct {
	Error    error
	Resource string
	Region   string
	Data     *schema.AWSResources
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
		prov, err := provider.CreateProvider(r.Builder.Config.Provider)
		if err != nil {
			return err
		}

		// Resources based
		for _, t := range r.Builder.Config.Resources {
			go func(name, region string) {
				re := Record{
					Error:    nil,
					Resource: name,
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

	result := schema.AWSResources{}
	for i := 0; i < r.TotalCount; i++ {
		record := <-ch

		if record.Data != nil {
			setReturnData(&result, record.Data, record.Resource)
		}

		if record.Error != nil {
			errors = append(errors, record.Error)
		}
	}

	if err := PrintScanResult(out, r.Builder.Config.Provider, result); err != nil {
		return err
	}
	return nil
}

// setReturnData sets return value to result
func setReturnData(result, data *schema.AWSResources, resource string) {
	switch resource {
	case "ec2":
		result.EC2 = append(result.EC2, data.EC2...)
	}
}

// PrintScanResult prints scan result
func PrintScanResult(out io.Writer, provider string, result schema.AWSResources) error {
	var scanData = struct {
		Summary  schema.AWSResources
		Provider string
	}{
		Summary:  result,
		Provider: provider,
	}

	funcMap := template.FuncMap{
		"decorate": color.DecorateAttr,
	}

	// Template for scan result
	w := tabwriter.NewWriter(out, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("Result").Funcs(funcMap).Parse(templates.Templates[provider]))

	err := t.Execute(w, scanData)
	if err != nil {
		return err
	}
	return w.Flush()
}
