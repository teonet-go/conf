// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Config helper go package. Number module parses numbers from strings to
// native values.

package conf

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

// Number is an interface that combines the Integer and Float interfaces from
// the constraints package. This allows code to accept numbers of either integer
// or floating point types through the single Number interface.
type Number interface {
	constraints.Integer | constraints.Float
}

// ParseNumber parses a string into a number of type T.
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
func ParseNumber[T Number](s string) (n T, err error) {

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

// NumberToString converts the input number to a string.
//
//	n: the input number of type T
//
// Returns the string representation of the input number.
func NumberToString[T Number](n T) string {
	return fmt.Sprintf("%v", n)
}
