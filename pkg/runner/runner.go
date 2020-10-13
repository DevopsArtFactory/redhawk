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

package runner

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/redhawk/pkg/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/client"
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/printer"
	"github.com/DevopsArtFactory/redhawk/pkg/provider"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
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
	return &Runner{}
}

// List retrieves resources in AWS
func (r Runner) List(out io.Writer) error {
	t := time.Now()
	logrus.Info("start scanning resources")

	var errors []error
	ch := make(chan Record)
	totalCount := 0
	for _, resource := range r.Builder.Config.Resources {
		if !resource.Global {
			totalCount += len(r.Builder.Config.Regions)
		} else {
			totalCount++
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
					Region:   constants.DefaultRegion,
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

	logrus.Debugf("Check format validation: %s", r.Builder.Flags.Output)
	if err := tools.CheckValidFormat(r.Builder.Flags.Output); err != nil {
		return err
	}

	logrus.Debugf("Create a printer for output: %s", r.Builder.Flags.Output)
	printer, err := printer.SelectPrinter(r.Builder.Flags.Output)
	if err != nil {
		return err
	}
	logrus.Debug("Printer is successfully created")

	logrus.Debugf("Set a number of data for printer: %d", len(result.Resources))
	pr, err := printer.SetData(result.Provider, result.Resources)
	if err != nil {
		return err
	}
	logrus.Debug("Data setting for printer is done")

	logrus.Debug("Start printer to print the result")
	if err := pr.Print(); err != nil {
		return err
	}

	if len(errors) > 0 && (logrus.GetLevel() == logrus.DebugLevel && logrus.GetLevel() == logrus.TraceLevel) {
		for _, err := range errors {
			logrus.Error(err.Error())
		}
	}

	end := time.Now()
	logrus.Infof("Scan time: %f sec", end.Sub(t).Seconds())

	return nil
}
