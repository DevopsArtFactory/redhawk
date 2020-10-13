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
	"html/template"
	"io"
	"os"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/redhawk/pkg/color"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/DevopsArtFactory/redhawk/pkg/templates"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
)

type StdOutPrinter struct {
	Out      io.Writer
	Provider string
	Data     map[string][]resource.Resource
}

// NewStdOutPrinter creates a new stdout printer
func NewStdOutPrinter() Printer {
	return StdOutPrinter{}
}

// SetData sets data
func (s StdOutPrinter) SetData(provider string, d []resource.Resource) (Printer, error) {
	ret := map[string][]resource.Resource{}
	for _, r := range d {
		rt := r.GetResource()
		if _, ok := ret[rt]; !ok {
			ret[rt] = []resource.Resource{}
		}

		ret[rt] = append(ret[rt], r)
	}
	s.Data = ret
	s.Out = os.Stdout
	s.Provider = provider
	return s, nil
}

// Print shows data to Standard Out
func (s StdOutPrinter) Print() error {
	detail := viper.GetBool("detail")
	logrus.Debugf("Detailed mode enabled: %t", detail)
	var scanData = struct {
		Summary  map[string][]resource.Resource
		Provider string
		Detail   bool
	}{
		Summary:  s.Data,
		Provider: s.Provider,
		Detail:   detail,
	}

	funcMap := template.FuncMap{
		"decorate": color.DecorateAttr,
		"format":   tools.Formatting,
	}

	// Template for scan result
	w := tabwriter.NewWriter(s.Out, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("Result").Funcs(funcMap).Parse(templates.Templates[s.Provider]))

	err := t.Execute(w, scanData)
	if err != nil {
		return err
	}
	return w.Flush()
}
