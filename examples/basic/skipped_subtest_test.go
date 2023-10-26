package basic

import (
	"testing"

	"github.com/bugVanisher/gunit"
)

func TestHowToSkipASubtest(t *testing.T) {
	gunit.Run(new(HowToSkipASubtestFixture), t)
}

type HowToSkipASubtestFixture struct {
	*gunit.Fixture
}

func (this *HowToSkipASubtestFixture) TestA() {
	this.SkipNow()
	this.Assert(false)
}
func (this *HowToSkipASubtestFixture) TestB() {}
func (this *HowToSkipASubtestFixture) TestC() {}
