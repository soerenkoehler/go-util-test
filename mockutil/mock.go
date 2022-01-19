package mockutil

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

type mockInvocation []interface{}

type mockWriter struct {
	writeMethod func(p []byte) (int, error)
}

func (m mockWriter) Write(p []byte) (int, error) {
	return m.writeMethod(p)
}

type argumentCheck func(interface{}) error

func createArgumentCheck(
	name string,
	expected interface{},
	transform func(interface{}) interface{}) argumentCheck {
	return func(actual interface{}) error {
		if transform(expected) == transform(actual) {
			return nil
		} else {
			return fmt.Errorf(
				"%s: expected = %v actual = %v",
				name,
				expected,
				actual)
		}
	}
}

func Any() argumentCheck {
	return createArgumentCheck(
		"Any",
		nil,
		func(v interface{}) interface{} {
			return nil
		})
}

func ArgValue(expected interface{}) argumentCheck {
	return createArgumentCheck(
		"Value",
		expected,
		func(v interface{}) interface{} {
			return v
		})
}

func ArgPointer(expected interface{}) argumentCheck {
	return createArgumentCheck(
		"Pointer",
		expected,
		func(v interface{}) interface{} {
			return reflect.ValueOf(v).Pointer()
		})
}

func ArgType(expected interface{}) argumentCheck {
	return createArgumentCheck(
		"Type",
		expected,
		func(v interface{}) interface{} {
			return reflect.TypeOf(v)
		})
}

type Registry struct {
	T           *testing.T
	invocations []mockInvocation
}

func (registry *Registry) Writer(name string) io.Writer {
	return mockWriter{
		writeMethod: func(p []byte) (int, error) {
			return len(p), nil
		}}
}

func (registry *Registry) Register(invocation ...interface{}) {
	registry.invocations = append(registry.invocations, invocation)
}

func (registry *Registry) NoMoreInvocations() *Registry {
	if remaining := len(registry.invocations); remaining != 0 {
		invocations := make([]string, 0, remaining)
		for _, invocation := range registry.invocations {
			invocations = append(invocations, fmt.Sprint(invocation))
		}
		registry.T.Errorf(
			"unexpected invocations:\n%v",
			strings.Join(invocations, "\n"))
	}

	return registry
}

func (registry *Registry) Verify(name string, args ...argumentCheck) *Registry {
	registered := registry.invocations[0]
	registry.invocations = registry.invocations[1:]

	if err := verifyInvocation(registered, name, args...); err != nil {
		registry.T.Errorf(
			"mismatch: %v\nregistered invocation: %v",
			err,
			registered)
	}

	return registry
}

func verifyInvocation(
	actualInvocation mockInvocation,
	expectedName string,
	expectedArgs ...argumentCheck) error {

	if expectedName != actualInvocation[0] {
		return fmt.Errorf("invocation name mismatch")
	}

	if len(actualInvocation)-1 != len(expectedArgs) {
		return fmt.Errorf("argument list length mismatch")
	}

	for index, expectedArg := range expectedArgs {
		if err := expectedArg(actualInvocation[1+index]); err != nil {
			return err
		}
	}

	return nil
}
