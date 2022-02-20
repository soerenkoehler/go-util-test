package mockutil

import (
	"fmt"
	"strings"
	"testing"
)

type mockInvocation []interface{}

type Registry struct {
	T           *testing.T
	StdOut      *strings.Builder
	StdErr      *strings.Builder
	invocations []mockInvocation
}

func NewRegistry(t *testing.T) Registry {
	return Registry{
		T:      t,
		StdOut: &strings.Builder{},
		StdErr: &strings.Builder{},
	}
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
		registry.T.Fatalf(
			"unexpected invocations:\n%v",
			strings.Join(invocations, "\n"))
	}

	return registry
}

func (registry *Registry) Verify(name string, args ...argumentCheck) *Registry {
	if len(registry.invocations) == 0 {
		registry.T.Fatalf(
			"missing invocation: %s %v",
			name,
			args)
	}

	registered := registry.invocations[0]
	registry.invocations = registry.invocations[1:]

	if err := verifyInvocation(registered, name, args...); err != nil {
		registry.T.Fatalf(
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
