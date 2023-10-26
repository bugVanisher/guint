package basic

import (
	"github.com/bugVanisher/gunit"
	"testing"
)

func TestMoreAssertFixture(t *testing.T) {
	gunit.Run(new(MoreAssertFixture), t)
}

type MoreAssertFixture struct {
	*gunit.Fixture
}

func (this *MoreAssertFixture) Setup() {
}

func (this *MoreAssertFixture) FixtureSetup() {
}

func (this *MoreAssertFixture) FixtureTeardown() {
}

func (this *MoreAssertFixture) Test() {
	// External assertion functions from the `should` package:
	// import "github.com/smarty/assertions/should"
	// go version should over 1.18
	// ...
	// this.So(42, should.Equal, 42)
	// this.So("Hello, World!", should.ContainSubstring, "World")
}
