package mock

import (
	"fmt"
	"reflect"
)

type argumentCheck func(interface{}) error

func createArgumentCheck(
	name string,
	expected interface{},
	transform func(interface{}) interface{},
	comparison func(x, y interface{}) bool) argumentCheck {

	return func(actual interface{}) error {
		transformedExpected := transform(expected)
		transformedActual := transform(actual)
		if comparison(transformedExpected, transformedActual) {
			return nil
		} else {
			return fmt.Errorf(
				"%s: expected = %v actual = %v",
				name,
				transformedExpected,
				transformedActual)
		}
	}
}

func equal(x, y interface{}) bool {
	return x == y
}

func Is(expected interface{}) argumentCheck {
	return createArgumentCheck(
		"Value",
		expected,
		func(v interface{}) interface{} { return v },
		reflect.DeepEqual)
}

func IsAny() argumentCheck {
	return createArgumentCheck(
		"Any",
		nil,
		func(v interface{}) interface{} { return nil },
		equal)
}

func IsFunc(expected interface{}) argumentCheck {
	return createArgumentCheck(
		"Func",
		expected,
		func(v interface{}) interface{} {
			value := reflect.ValueOf(v)
			if value.Kind() != reflect.Func {
				panic("expected func")
			}
			return value.Pointer()
		},
		equal)
}

func IsType(expected interface{}) argumentCheck {
	return createArgumentCheck(
		"Type",
		expected,
		func(v interface{}) interface{} { return reflect.TypeOf(v) },
		equal)
}
