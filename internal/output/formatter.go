package output

import "io"

type Formatter interface {
	Format(w io.Writer, data any) error
}

func New(format string) Formatter {
	switch format {
	case "table":
		return &TableFormatter{}
	case "csv":
		return &CSVFormatter{}
	default:
		return &JSONFormatter{}
	}
}
