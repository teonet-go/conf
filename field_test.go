package conf

import (
	"fmt"
	"testing"
)

// test for convertToNumber

func TestConvertToNumber(t *testing.T) {

	s := "1.32"
	f1, err := ParseNumber[float32](s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f1)

	s2 := "1.64"
	f2, err := ParseNumber[float64](s2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f2)

	type myFloat32 float32
	s3 := "1.32"
	f3, err := ParseNumber[float32](s3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(myFloat32(f3))

	type myFloat64 float64
	s4 := "1.64"
	f4, err := ParseNumber[myFloat64](s4)
	if err == nil {
		err = fmt.Errorf("myFloat64 is not a number, cust it to some Number format")
		t.Fatal(err)
	}
	fmt.Println(f4)

	s5 := "100"
	i1, err := ParseNumber[int](s5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i1)

	s6 := "257"
	// var i2 uint16
	i2, err := ParseNumber[uint8](s6)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i2)
}

func TestConvertToString(t *testing.T) {

	n := 3.14
	s := NumberToString(n)
	fmt.Println(s)

	type myInt int

	var i myInt = 75
	s2 := NumberToString(i)
	fmt.Println(s2)
}