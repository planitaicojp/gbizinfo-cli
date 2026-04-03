package output

import (
	"fmt"
	"reflect"
	"strings"
)

// flattenData takes a slice value and returns headers and rows.
// If an element struct contains exactly one slice-of-struct field,
// that field is expanded: each nested element becomes its own row
// with the parent's scalar fields repeated. Only the first such
// slice field is expanded; additional ones are rendered as strings.
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
		if isStructSlice(ft) && sliceIndex < 0 {
			sliceIndex = i
			sliceElemType = derefType(ft.Elem())
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
			scalarVals[j] = formatValue(row.Field(idx))
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
				record = append(record, formatValue(child.Field(f)))
			}
			rows = append(rows, record)
		}
	}

	return headers, rows
}

// isStructSlice returns true if t is a slice whose element type is
// a struct (or pointer to struct).
func isStructSlice(t reflect.Type) bool {
	if t.Kind() != reflect.Slice {
		return false
	}
	return derefType(t.Elem()).Kind() == reflect.Struct
}

// derefType follows a pointer type to its element type.
func derefType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

// formatValue converts a reflect.Value to a string suitable for
// table/CSV output. Slices of primitives are joined with "; ".
func formatValue(v reflect.Value) string {
	if v.Kind() == reflect.Slice && !isStructSlice(v.Type()) {
		parts := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			parts[i] = fmt.Sprintf("%v", v.Index(i).Interface())
		}
		return strings.Join(parts, "; ")
	}
	return fmt.Sprintf("%v", v.Interface())
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
