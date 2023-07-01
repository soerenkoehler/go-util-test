package test

import "testing"

type Suite map[string]func(*testing.T)

func RunSuite(
	t *testing.T,
	setup func(*testing.T),
	teardown func(*testing.T),
	suite Suite) {

	for testName, testFunc := range suite {
		if setup != nil {
			setup(t)
		}

		t.Run(testName, testFunc)

		if teardown != nil {
			teardown(t)
		}
	}
}
