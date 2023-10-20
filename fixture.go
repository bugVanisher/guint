// Package gunit provides "testing" package hooks and convenience
// functions for writing tests in an xUnit style.
// See the README file and the examples folder for examples.
package gunit

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/bugVanisher/gunit/reports"
)

// Fixture keeps track of test status (failed, passed, skipped) and
// handles custom logging for xUnit style tests as an embedded field.
// The Fixture itself and it's *testing.T control the test flow.
type Fixture struct {
	t               TestingT
	log             *bytes.Buffer
	verbose         bool
	testPackageName string
	logger          *Logger
	skipLeftTests   bool
	finished        bool
	parent          *Fixture
}

func newFixture(t TestingT, verbose bool, pkgName string) *Fixture {
	t.Helper()

	return &Fixture{t: t, verbose: verbose, log: &bytes.Buffer{}, testPackageName: pkgName, skipLeftTests: false}
}

// T exposes the TestingT (*testing.T) instance.
func (f *Fixture) T() TestingT { return f.t }

// Run is analogous to *testing.T.Run and allows for running subtests from
// test fixture methods (such as for table-driven tests).
func (f *Fixture) Run(name string, test func(fixture *Fixture)) {
	pkgName := retrieveTestPackageName()
	f.t.(*testing.T).Run(name, func(t *testing.T) {
		t.Helper()

		fixture := newFixture(t, f.verbose, pkgName)
		defer fixture.finalize()
		test(fixture)
	})
}

// So is a convenience method for reporting assertion failure messages,
// from the many assertion functions found in github.com/smarty/assertions/should.
// Example: this.So(actual, should.Equal, expected)
func (f *Fixture) So(actual interface{}, assert assertion, expected ...interface{}) bool {
	failure := assert(actual, expected...)
	failed := len(failure) > 0
	if failed {
		f.fail(failure)
	}
	return !failed
}

// Assert tests a boolean which, if not true, marks the current test case as failed and
// prints the provided message.
func (f *Fixture) Assert(condition bool, messages ...string) bool {
	if !condition {
		if len(messages) == 0 {
			messages = append(messages, "Expected condition to be true, was false instead.")
		}
		f.fail(strings.Join(messages, ", "))
	}
	return condition
}
func (f *Fixture) AssertEqual(expected, actual interface{}) bool {
	return f.Assert(expected == actual, fmt.Sprintf(comparisonFormat, fmt.Sprint(expected), fmt.Sprint(actual)))
}
func (f *Fixture) AssertSprintEqual(expected, actual interface{}) bool {
	return f.AssertEqual(fmt.Sprint(expected), fmt.Sprint(actual))
}
func (f *Fixture) AssertSprintfEqual(expected, actual interface{}, format string) bool {
	return f.AssertEqual(fmt.Sprintf(format, expected), fmt.Sprintf(format, actual))
}
func (f *Fixture) AssertDeepEqual(expected, actual interface{}) bool {
	return f.Assert(reflect.DeepEqual(expected, actual),
		fmt.Sprintf(comparisonFormat, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual)))
}

// Write implements io.Writer. There are rare times when this is convenient (debugging via `log.SetOutput(fixture)`).
func (f *Fixture) Write(p []byte) (int, error) { return f.log.Write(p) }
func (f *Fixture) Failed() bool                { return f.t.Failed() }
func (f *Fixture) Name() string                { return f.t.Name() }

func (f *Fixture) fail(failure string) {
	f.t.Fail()
	f.t.Log(reports.FailureReport(failure, reports.StackTrace()))
}

func (f *Fixture) finalize() {
	f.t.Helper()

	if r := recover(); r != nil {
		f.recoverPanic(r)
	}

	if f.t.Failed() || (f.verbose && f.log.Len() > 0) {
		f.t.Log("\n" + strings.TrimSpace(f.log.String()) + "\n")
	}
	f.finished = true
}
func (f *Fixture) recoverPanic(r interface{}) {
	f.t.Fail()
	f.t.Log(reports.PanicReport(r, debug.Stack()))
	f.skipLeftTests = true
}

func (f *Fixture) GetLogger() *Logger {
	return &Logger{t: f.t, testPackageName: f.testPackageName}
}

// FatalStop stop the test right now and the tests behind it will be skipped in SequentialTestCases mode.
func (f *Fixture) FatalStop(args ...interface{}) {
	f.skipLeftTests = true
	if f.finished {
		f.GetLogger().Error().Msg("[FatalStop]subtest has finished, parent test fail now")
		f.parent.skipLeftTests = true
		f.parent.T().(*testing.T).Fatal(args)
	}
	f.T().(*testing.T).Fatal(args...)
}

func (f *Fixture) FatalfStop(format string, args ...interface{}) {
	f.skipLeftTests = true
	if f.finished {
		f.GetLogger().Error().Msg("[FatalStop]subtest has finished, parent test fail now")
		f.parent.skipLeftTests = true
		f.parent.T().(*testing.T).Fatalf(format, args...)
	}
	f.T().(*testing.T).Fatalf(format, args...)
}

const comparisonFormat = "Expected: [%s]\nActual:   [%s]"

// assertion is a copy of github.com/smarty/assertions.assertion.
type assertion func(actual interface{}, expected ...interface{}) string
