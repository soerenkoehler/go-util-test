package mockutil

type mockInvocation []interface{}

type mockWriter struct {
	// lines []string
}

func (m mockWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
