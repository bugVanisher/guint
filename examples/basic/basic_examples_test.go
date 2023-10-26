package basic

import (
	"testing"
	"time"

	"github.com/bugVanisher/gunit"
)

func TestExampleFixture(t *testing.T) {
	gunit.Run(new(ExampleFixture), t, gunit.Options.SequentialTestCases())
}

type ExampleFixture struct {
	*gunit.Fixture // Required: Embedding this type is what allows gunit.Run to run the tests in this fixture.
	man            Person
	// Declare useful state here (probably the stuff being testing, any fakes, etc...).
}

type Person struct {
	Name  string
	Motto string
}

func (this *ExampleFixture) SetupStuff() {
	// This optional method will be executed before each "Test"
	// method because it starts with "Setup".
	this.GetLogger().Info().Str("before subtest", this.T().Name()).Msg("")
}
func (this *ExampleFixture) TeardownStuff() {
	// This optional method will be executed after each "Test"
	// method (because it starts with "Teardown"), even if the test method panics.
	this.GetLogger().Info().Str("after subtest", this.T().Name()).Msg("")
}

func (this *ExampleFixture) FixtureSetupStuff() {
	// This optional method will be executed before all "Test"
	// method because it starts with "FixtureSetup".
	this.GetLogger().Info().Str("before all subtest", this.T().Name()).Msg("")
	this.man = Person{
		Name:  "bugVanisher",
		Motto: "I want something just like this...",
	}
}
func (this *ExampleFixture) FixtureTeardownStuff() {
	// This optional method will be executed after all "Test"
	// method (because it starts with "FixtureTeardown"), even if the test method panics.
	this.GetLogger().Info().Str("after all subtest", this.T().Name()).Msg("")
}

// These are actual test cases:
func (this *ExampleFixture) TestWithLog() {
	this.GetLogger().Info().Any("man", this.man).Msg("")
	time.Sleep(1 * time.Second)
}

func (this *ExampleFixture) TestWithAssertions() {
	// Built-in assertion functions:
	this.Assert(true)
	this.AssertDeepEqual(1, 1)
	this.AssertSprintEqual(1, 1.0)
	this.AssertSprintfEqual(uint(1), int64(1), "%d")
}
