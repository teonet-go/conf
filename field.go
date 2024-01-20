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
	"unicode"

	"golang.org/x/exp/constraints"
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

// GetFields generates a list of fields from the given object and calls a
// function for each detected field.
//
// The function accepts two parameters:
//
//   - o: the object from which to extract the fields. It may be a struct or a
//     map[string]any.
//   - f: the function to be called for each field, which takes a pointer to a
//     Field[T] struct as its parameter. Where T is the type of the Entry fied
//     in the Field struct.
//
// It returns a Fields[T] which is a slice of pointers to Field[T] structs.
func GetFields[T any](o any, f func(field *Field[T])) (fields Fields[T]) {

	// Make description of fields
	makeField := func(fld reflect.Value, name, nameDisplay string) {
		fieldValue := fld.Interface()
		fieldType := fld.Type()
		field := &Field[T]{
			Name:        name,
			Value:       fieldValue,
			NameDisplay: nameDisplay,
			Type:        fieldType.String(),
			ValueStr:    fmt.Sprintf("%v", fieldValue),
		}
		fields = append(fields, field)
		f(field)

		fmt.Printf("%s: %v of type %v\n", field.Name, field.Value, field.Type)
	}

	// Make fields
	switch {

	// If the o object is struct
	case isStruct(o):
		v := reflect.ValueOf(o)
		t := reflect.TypeOf(o)
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name
			if field.CanInterface() {
				makeField(field, fieldName, fieldName)
			}
		}

	// If the o object is map
	case isMap(o):
		for key, val := range o.(map[string]interface{}) {
			v := reflect.ValueOf(val)
			t := reflect.TypeOf(val)
			fmt.Println("key: ", key, v, t)

			makeField(v, key, uppercaseFirstRune(key))
		}

	// If the o parameter is not a struct or a map, panic
	default:
		panic("o parameter of GetFields should be a struct or a map")
	}

	return
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
func SetValue[T any](p any, field *Field[T], value string) (err error) {

	// setError returns an error with the provided field name, value, and type.
	setError := func(name, value, t string) error {
		return fmt.Errorf("can't set %s: %v of type %s", name, value, t)
	}

	// setStructValue sets the value of a field in a struct.
	setStructValue := func(val reflect.Value) (err error) {
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
			// case "options.RadioGroup", "*options.RadioGroup":
			// 	// TODO: set real value here, if possible - use reflet
			// 	fmt.Println("radio group set value: ", value)
			// 	// val.Pointer()
			// 	val.Interface().(*options.RadioGroup).Selected = 33
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
				err = setError(field.Name, value, val.Type().String())
			}
		} else {
			err = setError(field.Name, value, val.Type().String())
		}
		return
	}

	// setMapValue sets the value of a field in a map.
	setMapValue := func(m map[string]any) (err error) {
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

	// Check if the p parameter is a pointer to a struct or a map, set values
	// for struct or map or panic if parameter is not valid (is not a pointer
	//to a struct or a map).
	switch {

	// If the p parameter is a pointer to a struct
	case isStructPtr(p):
		v := reflect.ValueOf(p).Elem()
		field := v.FieldByName(field.Name)
		err = setStructValue(field)

	// If the p parameter is a struct
	case isMap(p) || isMapPtr(p):
		var m map[string]any
		if isMapPtr(p) {
			m = *p.(*map[string]any)
		} else {
			m = p.(map[string]any)
		}
		err = setMapValue(m)

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
func ValidateValue[T any](field *Field[T], value string) (err error) {

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

// uppercaseFirstRune converts the first rune of a string to uppercase.
//
// It takes a string as a parameter and returns the modified string.
func uppercaseFirstRune(s string) string {
	runes := []rune(s)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

// Fields is a slice of Field pointers for the generic type T.
// It allows operating on a group of fields together.
type Fields[T any] []*Field[T]

// SetValues iterates over each field in the Fields collection and sets their values
// based on the provided function.
//
// Parameters:
//   - p: The target object where the field values will be set.
//   - f: A function that takes a pointer to a Field and returns a string representing the field value.
func (fields Fields[T]) SetValues(p any, f func(field *Field[T]) string) {
	for _, field := range fields {
		txt := f(field)
		SetValue(p, field, txt)
	}
}

// Number is an interface that combines the Integer and Float interfaces from
// the constraints package. This allows code to accept numbers of either integer
// or floating point types through the single Number interface.
type Number interface {
	constraints.Integer | constraints.Float
}

// parseNumber parses a string into a number of type T.
//
// The function takes a string `s` as input and attempts to convert it into a
// number of type T. The function supports number types such as int, int8,
// int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, and
// float64.
//
// Parameters:
//   - s: the string to be parsed.
//
// Returns:
//   - n: the parsed number of type T.
//   - err: an error indicating any parsing failure.
func parseNumber[T Number](s string) (n T, err error) {

	parseFloat := func(s string, bitsize int) (T, error) {
		var num float64
		num, err = strconv.ParseFloat(s, bitsize)
		if err != nil {
			return T(0), err
		}
		return T(num), nil
	}

	switch any(n).(type) {

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, err := strconv.Atoi(s)
		return T(i), err

	case float32:
		return parseFloat(s, 32)

	case float64:
		return parseFloat(s, 64)

	default:
		err = fmt.Errorf("can't convert type %T to number", n)
	}

	return
}
