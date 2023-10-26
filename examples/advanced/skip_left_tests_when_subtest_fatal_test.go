package advanced

import (
	"github.com/bugVanisher/gunit"
	"testing"
)

//TestSubtestFatal other subtest won't be skipped when subtest fatal in parallel mode.
func TestSubtestFatal(t *testing.T) {
	gunit.Run(new(SubtestFatal), t)
}

func TestSubtestFatalSequential(t *testing.T) {
	gunit.Run(new(SubtestFatal), t, gunit.Options.SequentialTestCases())
}

type SubtestFatal struct {
	*gunit.Fixture
}

func (f *SubtestFatal) Setup() {
}

func (f *SubtestFatal) FixtureSetup() {
}

func (f *SubtestFatal) FixtureTeardown() {
}

func (f *SubtestFatal) Test1() {
}
func (f *SubtestFatal) Test2() {
	f.FatalStop()
}

//Test3 will only be skipped in Sequential mode
func (f *SubtestFatal) Test3() {
}

//Test4 will only be skipped in Sequential mode
func (f *SubtestFatal) Test4() {
}
