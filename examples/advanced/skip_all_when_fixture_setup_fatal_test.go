package advanced

import (
	"github.com/bugVanisher/gunit"
	"testing"
)

func TestFixtureSetupFail(t *testing.T) {
	gunit.Run(new(FixtureSetupFail), t)
}

func TestFixtureSetupFailSequential(t *testing.T) {
	gunit.Run(new(FixtureSetupFail), t, gunit.Options.SequentialTestCases())
}

type FixtureSetupFail struct {
	*gunit.Fixture
}

func (f *FixtureSetupFail) Setup() {
}

func (f *FixtureSetupFail) FixtureSetup() {
	f.FatalStop()
}

func (f *FixtureSetupFail) FixtureTeardown() {
}

//Test1 will be skipped
func (f *FixtureSetupFail) Test1() {
}

//Test2 will be skipped
func (f *FixtureSetupFail) Test2() {
}
