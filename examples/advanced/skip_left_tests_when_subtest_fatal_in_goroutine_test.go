package advanced

import (
	"github.com/bugVanisher/gunit"
	"testing"
	"time"
)

//TestSubtestFatalInGoroutineSequential other subtest will be skipped when subtest fatal in goroutine in Sequential mode.
func TestSubtestFatalInGoroutineSequential(t *testing.T) {
	gunit.Run(new(SubtestFatalInGoroutine), t, gunit.Options.SequentialTestCases())
}

type SubtestFatalInGoroutine struct {
	ExtendFixture
}

func (f *SubtestFatalInGoroutine) Setup() {
}

func (f *SubtestFatalInGoroutine) FixtureSetup() {
}

func (f *SubtestFatalInGoroutine) FixtureTeardown() {
}

func (f *SubtestFatalInGoroutine) Test1() {
}
func (f *SubtestFatalInGoroutine) Test2() {
	f.Streaming()
}

//Test3 will only be skipped in Sequential mode
func (f *SubtestFatalInGoroutine) Test3() {
	time.Sleep(500 * time.Millisecond)
}

//Test4 will only be skipped in Sequential mode
func (f *SubtestFatalInGoroutine) Test4() {
}

//Test4 will only be skipped in Sequential mode
func (f *SubtestFatalInGoroutine) Test5() {
}

//Test4 will only be skipped in Sequential mode
func (f *SubtestFatalInGoroutine) Test6() {
}
