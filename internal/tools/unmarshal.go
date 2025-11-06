// Package tools provides utility functions
package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type RequiredTag struct {
	Struct any
}

func NewRequiredTagStruct(strct any) *RequiredTag {
	return &RequiredTag{
		Struct: strct,
	}
}

func (r *RequiredTag) UnmarhsalJSON(data []byte) error {
	// Unmarshal into the target struct first
	if err := json.Unmarshal(data, r.Struct); err != nil {
		return err
	}

	// Perform validation based on "required" tag
	val := reflect.ValueOf(r.Struct).Elem() // Get the underlying struct value
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		_, ok := field.Tag.Lookup("required")
		if ok {
			// Check if the field is a zero value (e.g., empty string, zero int, nil pointer)
			if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(field.Type).Interface()) {
				return fmt.Errorf("field '%s' is required but is empty", field.Name)
			}
		}
	}
	return nil
}

func (r *RequiredTag) MarshalJSON() ([]byte, error) {
	// Perform validation based on "required" tag
	val := reflect.ValueOf(r.Struct).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		_, ok := field.Tag.Lookup("required")
		if ok {
			if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(field.Type).Interface()) {
				return nil, fmt.Errorf("field '%s' is required but is empty", field.Name)
			}
		}
	}
	// Marshal the source struct if validation passes
	return json.Marshal(r.Struct)
}
