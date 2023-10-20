package gunit

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/bugVanisher/gunit/scan"
)

// Run receives an instance of a struct that embeds *Fixture.
// The struct definition may include Setup*, Teardown*, FixtureSetup*, FixtureTeardown* and Test*
// methods which will be run as an xUnit-style test fixture.
func Run(fixture interface{}, t *testing.T, options ...option) {
	t.Helper()

	if strings.Contains(runtime.Version(), "go1.14") {
		options = allSequentialForGo1Dot14(options)
	}
	pkgName := retrieveTestPackageName()
	run(fixture, t, newConfig(options...), pkgName)
}

func allSequentialForGo1Dot14(options []option) []option {
	// HACK to accommodate for https://github.com/bugVanisher/gunit/issues/28
	// Also see: https://github.com/golang/go/issues/38050
	return append(options, Options.AllSequential())
}

func run(fixture interface{}, t *testing.T, config configuration, pkgName string) {
	t.Helper()

	ensureEmbeddedFixture(fixture, t)

	_, filename, _, _ := runtime.Caller(2)
	positions := scan.LocateTestCases(filename)

	runner := newFixtureRunner(fixture, t, config, positions, pkgName)
	runner.ScanFixtureForTestCases()
	runner.RunTestCases()
}

func ensureEmbeddedFixture(fixture interface{}, t TestingT) {
	fixtureType := reflect.TypeOf(fixture)
	embedded, _ := fixtureType.Elem().FieldByName("Fixture")
	if embedded.Type != embeddedGoodExample.Type {
		t.(*testing.T).Fatalf("Type (%v) lacks embedded *gunit.Fixture.", fixtureType)
	}
}

var embeddedGoodExample, _ = reflect.TypeOf(new(struct{ *Fixture })).Elem().FieldByName("Fixture")
