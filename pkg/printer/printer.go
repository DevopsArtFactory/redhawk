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
	"fmt"

	"github.com/DevopsArtFactory/redhawk/pkg/resource"
)

type Printer interface {
	Print() error
	SetData(string, []resource.Resource) (Printer, error)
}

// SelectPrinter creates new printers
func SelectPrinter(outputType string) (Printer, error) {
	f, ok := printers[outputType]
	if !ok {
		return nil, fmt.Errorf("output type %s is not available for printer", outputType)
	}
	return f(), nil
}
