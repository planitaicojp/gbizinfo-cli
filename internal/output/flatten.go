package output

import (
	"fmt"
	"reflect"
	"strings"
)

// flattenData takes a slice value and returns headers and rows.
// Nested slice-of-struct fields are expanded: each nested element
// becomes its own row with the parent's scalar fields repeated.
func flattenData(val reflect.Value) (headers []string, rows [][]string) {
	if val.Len() == 0 {
		return nil, nil
	}

	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()

	var scalarIndices []int
	sliceIndex := -1
	var sliceElemType reflect.Type

	for i := 0; i < elemType.NumField(); i++ {
		ft := elemType.Field(i).Type
		if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Struct {
			sliceIndex = i
			sliceElemType = ft.Elem()
		} else {
			scalarIndices = append(scalarIndices, i)
		}
	}

	for _, idx := range scalarIndices {
		headers = append(headers, jsonTagName(elemType.Field(idx)))
	}
	if sliceIndex >= 0 {
		for i := 0; i < sliceElemType.NumField(); i++ {
			headers = append(headers, jsonTagName(sliceElemType.Field(i)))
		}
	}

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}

		scalarVals := make([]string, len(scalarIndices))
		for j, idx := range scalarIndices {
			scalarVals[j] = fmt.Sprintf("%v", row.Field(idx).Interface())
		}

		if sliceIndex < 0 {
			rows = append(rows, scalarVals)
			continue
		}

		nested := row.Field(sliceIndex)
		if nested.Len() == 0 {
			record := make([]string, len(headers))
			copy(record, scalarVals)
			rows = append(rows, record)
			continue
		}

		for k := 0; k < nested.Len(); k++ {
			child := nested.Index(k)
			if child.Kind() == reflect.Ptr {
				child = child.Elem()
			}
			record := make([]string, 0, len(headers))
			record = append(record, scalarVals...)
			for f := 0; f < child.NumField(); f++ {
				record = append(record, fmt.Sprintf("%v", child.Field(f).Interface()))
			}
			rows = append(rows, record)
		}
	}

	return headers, rows
}

func jsonTagName(field reflect.StructField) string {
	name := field.Tag.Get("json")
	if idx := strings.Index(name, ","); idx != -1 {
		name = name[:idx]
	}
	if name == "" || name == "-" {
		name = field.Name
	}
	return name
}
