package uiapi

import (
	"errors"
	"fmt"
	"strings"
)

func (it *ModelSpec) Valid() error {
	if len(it.Fields) == 0 {
		return errors.New("fields not found")
	}
	names := map[string]bool{}
	for _, field := range it.Fields {
		if field == nil {
			return errors.New("field empty")
		}
		name := strings.ToLower(field.Name)
		if _, ok := names[name]; ok {
			return errors.New("field name (" + field.Name + ") conflicted")
		}
		if field.Type == 0 {
			field.Type = ModelSpec_STRING
		}
		names[name] = true
	}
	for _, row := range it.DefaultRows {
		if row == nil {
			return errors.New("row empty")
		}

		if len(row.Values) != len(it.Fields) {
			return fmt.Errorf("row len(values %d) != len(fields %d)", len(row.Values), len(it.Fields))
		}
	}
	return nil
}
