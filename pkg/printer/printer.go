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
