package datax

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

func NewModelSpec(name string) *ModelSpec {
	return &ModelSpec{
		Name: name,
	}
}

func (it *ModelSpec) NewField(name string, args ...interface{}) *ModelSpec_Field {
	var field *ModelSpec_Field
	for _, f := range it.Fields {
		if name == f.Name {
			field = f
			break
		}
	}
	if field == nil {
		field = &ModelSpec_Field{
			Name: name,
		}
		it.Fields = append(it.Fields, field)
	}
	for _, arg := range args {
		switch arg.(type) {
		case ModelSpec_FieldType:
			field.Type = arg.(ModelSpec_FieldType)
		}
	}
	return field
}

func (it *ModelSpec_Field) SetType(v ModelSpec_FieldType) *ModelSpec_Field {
	it.Type = v
	return it
}

func (it *ModelSpec_Field) SetDefaultValue(v string) *ModelSpec_Field {
	it.DefaultValue = v
	return it
}

func (it *ModelSpec_Field) SetLength(v uint32) *ModelSpec_Field {
	it.Length = v
	return it
}

func (it *ModelSpec_Field) SetDecimalSize(precision, scale uint32) *ModelSpec_Field {
	it.DecimalSize = []uint32{precision, scale}
	return it
}

func (it *ModelSpec_Field) SetTermNs(v string) *ModelSpec_Field {
	it.TermNs = v
	return it
}

func (it *ModelSpec) NewRow(fields map[string]interface{}) *ModelSpec_Row {
	values := make([]string, len(it.Fields))
	for name, value := range fields {
		for i, field := range it.Fields {
			if name == field.Name {
				values[i] = fmt.Sprintf("%v", value)
			}
		}
	}
	return &ModelSpec_Row{
		Values: values,
	}
}
