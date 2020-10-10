package printer

var (
	printers = map[string]func() Printer{
		"stdout": NewStdOutPrinter,
		"csv":    NewCSVPrinter,
	}
)
