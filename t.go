package gunit

// TestingT represents the functional subset from *testing.T needed by Fixture.
type TestingT interface {
	Helper()
	Name() string
	Log(args ...interface{})
	Failed() bool
}
