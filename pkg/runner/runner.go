package runner

import (
	"github.com/DevopsArtFactory/redhawk/pkg/printer"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"

	"github.com/DevopsArtFactory/redhawk/pkg/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/client"
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/provider"
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
	Data     []resource.Resource
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
	totalCount := 0
	for _, resource := range r.Builder.Config.Resources {
		if !resource.Global {
			totalCount += len(r.Builder.Config.Regions)
		} else {
			totalCount += 1
		}
	}
	r.TotalCount = totalCount

	logrus.Debugf("Resource count is %d", totalCount)

	// Resources based
	for _, t := range r.Builder.Config.Resources {
		// Create new provider
		prov, err := provider.CreateProvider(r.Builder.Config.Provider)

		if err != nil {
			return err
		}

		if t.Global {
			logrus.Debugf("scanning global resources: %s", t.Name)
			go func(name string) {
				re := Record{
					Error:    nil,
					Resource: name,
				}
				c, err := prov.CreateClient(constants.DefaultRegion, name)
				if err != nil {
					re.Error = err
				} else {
					data, err := c.Scan()
					re.Error = err
					re.Data = data
				}
				ch <- re
			}(t.Name)
		} else {
			// Region based
			for _, region := range r.Builder.Config.Regions {
				logrus.Debugf("scanning regional resources: %s / %s", region, t.Name)
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
	}

	result := resource.Resources{
		Provider: r.Builder.Config.Provider,
	}

	for i := 0; i < r.TotalCount; i++ {
		record := <-ch

		if record.Data != nil {
			logrus.Debugf("data found: %s / %s / %s", record.Region, record.Resource, record.Data[0].GetResource())
			result.Resources = append(result.Resources, record.Data...)
		}

		if record.Error != nil {
			errors = append(errors, record.Error)
		}
	}

	logrus.Debugf("Completed gathering all data")

	outputFormat := viper.GetString("output")
	if err := tools.CheckValidFormat(outputFormat); err != nil {
		return err
	}

	printer, err := printer.SelectPrinter(outputFormat)
	if err != nil {
		return err
	}

	pr, err := printer.SetData(result.Provider, result.Resources)
	if err != nil {
		return err
	}

	if err := pr.Print(); err != nil {
		return err
	}

	return nil
}
