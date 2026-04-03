package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, data any) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		return fmt.Errorf("CSVフォーマットにはスライスが必要です")
	}
	if val.Len() == 0 {
		return nil
	}

	headers, rows := flattenData(val)

	writer := csv.NewWriter(w)
	if err := writer.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
