package advanced

import (
	"github.com/bugVanisher/gunit"
	"testing"
)

func TestExtendFixture(t *testing.T) {
	gunit.Run(new(ExtendFixture), t)
}

type ExtendFixture struct {
	FancyFixture
}

func (f *ExtendFixture) Setup() {
}

func (f *ExtendFixture) FixtureSetup() {
	f.StartLiving()
	f.GetMysqlClient().Exec("select 1;")
}

func (f *ExtendFixture) FixtureTeardown() {
	f.Release()
}

func (f *ExtendFixture) Test1() {
}

func (f *ExtendFixture) Test2() {
}
