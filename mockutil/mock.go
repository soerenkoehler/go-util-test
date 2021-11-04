package mockutil

import (
	"fmt"
	"reflect"
	"testing"
)

// Call is an array of type interface{} which describes a single mocked
// function call.
type Call []interface{}

type ArgumentCheck func(interface{}) error

// Registry collects a series of mocked calls.
type Registry struct {
	T     *testing.T // TODO how?
	calls []Call
}

// Register adds a MockCall to the array of registered calls.
func (registry *Registry) Register(call ...interface{}) {
	registry.calls = append(registry.calls, call)
}

func Any() ArgumentCheck {
	return func(actual interface{}) error {
		return nil
	}
}

func Is(expected interface{}) ArgumentCheck {
	return func(actual interface{}) error {
		if expected == actual {
			return nil
		} else {
			return fmt.Errorf(
				"expected value %v but got %v",
				expected,
				actual)
		}
	}
}

func SameType(expected interface{}) ArgumentCheck {
	return func(actual interface{}) error {
		expectedType := reflect.TypeOf(expected)
		actualType := reflect.TypeOf(actual)
		if expectedType == actualType {
			return nil
		} else {
			return fmt.Errorf(
				"expected type %v but got %v",
				expectedType,
				actualType)
		}
	}
}

// Verify checks if 'expectedCall' is registered at the first position of
// the call array.
func (registry *Registry) Verify(name string, args ...ArgumentCheck) *Registry {
	registeredCall := registry.calls[0]
	registry.calls = registry.calls[1:]

	if err := verifyCall(registeredCall, name, args...); err != nil {
		registry.T.Errorf(
			"mismatch: %v\nregistered call: %v",
			err,
			registeredCall)
	}

	return registry
}

func verifyCall(registeredCall Call, name string, args ...ArgumentCheck) error {

	if len(registeredCall)-1 != len(args) {
		return fmt.Errorf("argument list length mismatch")
	}

	for index, expectedArg := range args {
		if err := expectedArg(registeredCall[1+index]); err != nil {
			return err
		}
	}

	return nil
}
