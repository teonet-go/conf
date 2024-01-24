package conf

import (
	"fmt"
	"reflect"
	"unicode"
)

// Fields is a slice of Field pointers for the generic type T.
// It allows operating on a group of fields together.
type Fields[T any] []*Field[T]

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

// SetValues iterates over each field in the Fields collection and sets their
// values based on the provided function.
//
// Parameters:
//   - p: The target object where the field values will be set.
//   - f: A function that takes a pointer to a Field and returns a string
//     representing the field value.
func (fields Fields[T]) SetValues(p any, f func(field *Field[T]) (string, bool)) {
	for _, field := range fields {
		if txt, isStr := f(field); isStr {
			field.SetValue(p, txt)
			continue
		}
		field.SetValue(p)
	}
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
