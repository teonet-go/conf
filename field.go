package main

import (
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

type Field[T any] struct {
	Name     string
	Type     string
	ValueStr string
	Value    any
	Entry    T
}

// GetFields get struct fields from input p struct.

// GetFields generates a list of fields from the given object and calls a
// function for each detected field.
//
// The function accepts two parameters:
//   - o: the object from which to extract the fields.
//   - f: the function to be called for each field, which takes a pointer to a
//     Field[T] struct as its parameter.
//
// It returns a slice of pointers to Field[T] structs.
func GetFields[T any](o any, f func(field *Field[T])) (fields Fields[T]) {
	v := reflect.ValueOf(o)
	t := reflect.TypeOf(o)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		if field.CanInterface() {
			fieldValue := field.Interface()
			fieldType := field.Type()
			field := &Field[T]{
				Name: fieldName,
				Type: fieldType.String(),
				// TypeOf:   fieldType,
				ValueStr: fmt.Sprintf("%v", fieldValue),
				Value:    fieldValue,
			}
			fields = append(fields, field)
			f(field)

			fmt.Printf("%s: %v of type %v\n", fieldName, fieldValue, fieldType)
		}
	}
	return
}

func SetValue(p any, name string, value string) {

	// Get the reflect.Value of the p variable
	v := reflect.ValueOf(p).Elem()

	// Get the reflect.Value of the field by its name
	field := v.FieldByName(name)

	// Check if the field is valid and can be set. And set the new value for
	// the field by type.
	if field.IsValid() && field.CanSet() {
		switch field.Type().String() {
		case "string":
			field.SetString(value)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8",
			"uint16", "uint32", "uint64":
			i, _ := strconv.Atoi(value)
			field.SetInt(int64(i))
		case "float32":
			f, _ := strconv.ParseFloat(value, 32)
			field.SetFloat(f)
		case "float64":
			f, _ := strconv.ParseFloat(value, 64)
			field.SetFloat(f)
		default:
			fmt.Printf("can't set %s: %v of type %s\n", name, value,
				field.Type())
		}
	}
}

type Fields[T any] []*Field[T]

func (fields Fields[T]) SetValues(p any, f func(field *Field[T]) string) {
	for _, field := range fields {
		txt := f(field)
		SetValue(p, field.Name, txt)
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

// ParseNumber converts string to number.
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
