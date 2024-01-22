// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Config helper go package. It provides configuration handling functionality.
package conf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Field is a struct that contains metadata and values for a single field of a
// struct or map[string]any that is passed to the GetFields function. It
// allows associating additional data with each field through the generic
// Entry field. The NameDisplay field contains a display name for the field
// that can be used in UIs.
type Field[T any] struct {
	NameDisplay string // Name to show in form etc.
	Name        string // Field name
	Type        string // Field type
	ValueStr    string // Field value as string
	Value       any    // Field real value

	// Field entry is a custom field which can be used in GetFields and
	//SetValues callbacks
	Entry T
}

// SetValue sets the value of a field in a struct or map.
//
// The function takes in the following parameters:
//
//   - p: An interface{} value representing the struct or map.
//   - field: A pointer to the Field[T] struct, which contains information about
//     the field to be set.
//   - value: A string representing the value to be set.
//
// The function returns an error if the field cannot be set. The error message
// provides information about the field name, value, and type.
//
// The function supports setting values for fields of the following types:
// string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
// float32, and float64.
//
// If the value cannot be converted to the field's type, an error is returned.
//
// If the parameter 'p' is not a pointer to a struct or a map, the function
// panics.
//
// The function does not modify the 'p' parameter directly, but it modifies the
// value of the specified field.
//
// The function does not have any return values.
func (field *Field[T]) SetValue(p any, value string) (err error) {

	// Check if the p parameter is a pointer to a struct or a map, set values
	// for struct or map or panic if parameter is not valid (is not a pointer
	//to a struct or a map).
	switch {

	// If the p parameter is a pointer to a struct than set struct values
	case isStructPtr(p):
		err = field.setStructValue(p, value)

	// If the p parameter is map or pointer to map than set map values
	case isMap(p) || isMapPtr(p):
		var m map[string]any
		if isMapPtr(p) {
			m = *p.(*map[string]any)
		} else {
			m = p.(map[string]any)
		}
		err = field.setMapValue(m, value)

	// If the p parameter is not a pointer to a struct or a map, panic
	default:
		panic("p parameter of SetValue should be a pointer to struct or a map")
	}

	return
}

// ValidateValue validates the value of a given field.
//
// The function takes in two parameters: field, a pointer to a Field[T] struct,
// and value, a string representing the value to be validated. The field struct
// contains information about the field's type and name. The value parameter is
// a string that will be converted to the appropriate type based on the field's
// type.
//
// The function returns an error if the value is not of the expected type.
func (field *Field[T]) ValidateValue(value string) (err error) {

	switch field.Type {

	// Check integer
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8",
		"uint16", "uint32", "uint64":
		_, err = strconv.Atoi(value)

	// Check float
	case "float32", "float64":
		_, err = strconv.ParseFloat(value, 64)
	}

	// Check error
	if err != nil {
		err = fmt.Errorf("type of %s value should be %s", field.Name, field.Type)
	}

	return
}

// setStructValue sets the value of a struct field from string value.
func (field *Field[T]) setStructValue(p any, value string) (err error) {

	v := reflect.ValueOf(p).Elem() // Struct value
	name := field.Name             // Struct field name
	val := v.FieldByName(name)     // Struct field value

	if val.IsValid() && val.CanSet() {
		switch val.Type().String() {
		case "string":
			val.SetString(value)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8",
			"uint16", "uint32", "uint64":
			i, _ := strconv.Atoi(value)
			val.SetInt(int64(i))
		case "float32":
			f, _ := strconv.ParseFloat(value, 32)
			val.SetFloat(f)
		case "float64":
			f, _ := strconv.ParseFloat(value, 64)
			val.SetFloat(f)
		case "bool":
			val.SetBool(value == "true")
		case "[]int":
			s := strings.Split(strings.Trim(value, "[]"), " ")
			a := make([]int, len(s))
			for i, str := range s {
				val, _ := strconv.Atoi(str)
				a[i] = val
			}
			val.Set(reflect.ValueOf(a))
		case "[]float64":
			s := strings.Split(strings.Trim(value, "[]"), " ")
			a := make([]float64, len(s))
			for i, str := range s {
				val, _ := strconv.ParseFloat(str, 64)
				a[i] = val
			}
			val.Set(reflect.ValueOf(a))
		default:
			err = setError(name, value, val.Type().String())
		}
	} else {
		err = setError(name, value, val.Type().String())
	}
	return
}

// setMapValue sets the map field value from string value.
func (field *Field[T]) setMapValue(m map[string]any, value string) (err error) {
	key := field.Name

	switch field.Type {
	case "string":
		m[key] = value
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8",
		"uint16", "uint32", "uint64":
		i, _ := strconv.Atoi(value)
		m[key] = int64(i)
	case "float32":
		f, _ := strconv.ParseFloat(value, 32)
		m[key] = f
	case "float64":
		f, _ := strconv.ParseFloat(value, 64)
		m[key] = f
	case "bool":
		m[key] = value == "true"
	case "[]interface {}":
		s := strings.Split(strings.Trim(value, "[]"), " ")
		a := make([]any, len(s))
		for i, str := range s {
			val, err := strconv.Atoi(str)
			if err != nil {
				val, _ := strconv.ParseFloat(str, 64)
				a[i] = val
				continue
			}
			a[i] = val
		}
		m[key] = a
	default:
		err = setError(field.Name, value, field.Type)
	}
	return
}

// setError returns an error with the provided field name, value, and type.
func setError(name, value, t string) error {
	return fmt.Errorf("can't set %s: %v of type %s", name, value, t)
}

// isStruct checks if the given object is a struct.
//
// o: the object to check
//
// It returns true if the object is a struct, false otherwise.
func isStruct(o any) bool {
	return reflect.TypeOf(o).Kind() == reflect.Struct
}

// isStructPtr checks if the given object is a pointer to a struct.
//
// o: the object to check
//
// It returns true if the object is a pointer to a struct, false otherwise.
func isStructPtr(o any) bool {
	return reflect.TypeOf(o).Kind() == reflect.Ptr &&
		reflect.TypeOf(o).Elem().Kind() == reflect.Struct
}

// isMap checks if the given parameter is a map[string]interface{}.
//
// o: the parameter to be checked.
//
// It returns a boolean value indicating whether o is a map[string]interface{}.
func isMap(o any) bool {
	_, ok := o.(map[string]any)
	return ok
}

// isMapPtr checks if the given value is a pointer to a map.
//
//	o - the value to be checked.
//
//	Returns true if the value is a pointer to a map, false otherwise.
func isMapPtr(o any) bool {
	_, ok := o.(*map[string]any)
	return ok
}
