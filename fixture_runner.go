package gunit

import (
	"os"
	"os/signal"
	"reflect"
	"testing"

	"github.com/bugVanisher/gunit/scan"
)

const FixtureParallel = "FixtureParallel"

func newFixtureRunner(
	fixture interface{},
	outerT *testing.T,
	config configuration,
	positions scan.TestCasePositions,
	pkgName string,
) *fixtureRunner {
	outerT.Parallel()
	return &fixtureRunner{
		config:          config,
		fixtureSetup:    -1,
		fixtureTeardown: -1,
		setup:           -1,
		teardown:        -1,
		outerT:          outerT,
		fixtureType:     reflect.ValueOf(fixture).Type(),
		fixture:         reflect.New(reflect.ValueOf(fixture).Type().Elem()),
		positions:       positions,
		packageName:     pkgName,
	}
}

type fixtureRunner struct {
	outerT      *testing.T
	fixtureType reflect.Type
	fixture     reflect.Value

	fixtureSetup    int
	fixtureTeardown int
	config          configuration
	setup           int
	teardown        int
	focus           []*testCase
	tests           []*testCase
	positions       scan.TestCasePositions
	packageName     string
	skipLeftTests   bool
	skipChannel     chan bool
}

func (this *fixtureRunner) ScanFixtureForTestCases() {
	for methodIndex := 0; methodIndex < this.fixtureType.NumMethod(); methodIndex++ {
		methodName := this.fixtureType.Method(methodIndex).Name
		this.scanFixtureMethod(methodIndex, this.newFixtureMethodInfo(methodName))
	}
}

func (this *fixtureRunner) scanFixtureMethod(methodIndex int, method fixtureMethodInfo) {
	switch {
	case method.isFixTureSetup:
		this.fixtureSetup = methodIndex
	case method.isFixTureTeardown:
		this.fixtureTeardown = methodIndex
	case method.isSetup:
		this.setup = methodIndex
	case method.isTeardown:
		this.teardown = methodIndex
	case method.isFocusTest:
		this.focus = append(this.focus, this.buildTestCase(methodIndex, method))
	case method.isTest:
		this.tests = append(this.tests, this.buildTestCase(methodIndex, method))
	}
}

func (this *fixtureRunner) buildTestCase(methodIndex int, method fixtureMethodInfo) *testCase {
	return newTestCase(methodIndex, method, this.config, this.positions)
}

func (this *fixtureRunner) skipLeftTestsLoop() {
	for {
		select {
		case skip := <-this.skipChannel:
			this.skipLeftTests = skip
		default:

		}
	}
}

func (this *fixtureRunner) RunTestCases() {
	this.outerT.Helper()

	// Init Fixture for fixtureSetup and fixtureTeardown
	this.skipChannel = make(chan bool)
	defer close(this.skipChannel)
	go this.skipLeftTestsLoop()
	tmpFixture := newFixture(this.outerT, testing.Verbose(), this.packageName, this.skipChannel)
	defer this.leftTestsHandler()
	this.setInnerFixture(tmpFixture)
	defer tmpFixture.finalize()
	defer this.runFixtureTeardown()
	// Start goroutine to listen for SIGINT signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		// Clean up and exit
		this.runFixtureTeardown()
		this.outerT.Error("Interrupted by user")
		os.Exit(1)
	}()
	this.runFixtureSetup()

	if len(this.focus) > 0 {
		this.tests = append(this.focus, skipped(this.tests)...)
	}

	if this.config.SkippedTestCases {
		for _, test := range this.tests {
			this.outerT.Run(test.description, test.skip)
		}
	} else {
		if len(this.tests) > 0 {
			this.runTestCases()
		} else {
			this.outerT.Skipf("Fixture (%v) has no test cases.", this.fixtureType)
		}
	}
	// fix: replace inner testing
	this.setInnerFixture(tmpFixture)

}

func (this *fixtureRunner) runTestCases() {
	this.outerT.Helper()
	runCases := func(t *testing.T) {
		lastRunIndex := 0
		for _, test := range this.tests {
			lastRunIndex++
			if this.skipLeftTests {
				t.Run(test.description, test.skip)
				continue
			}
			test.Prepare(this.setup, this.teardown, this.fixture, this.packageName, this.skipChannel)
			test.Run(t)
		}
		this.tests = this.tests[lastRunIndex:]
	}
	if this.config.ParallelTestCases() {
		this.outerT.Run(FixtureParallel, func(innerT *testing.T) {
			runCases(innerT)
		})
	} else {
		runCases(this.outerT)
	}
}

func skipped(cases []*testCase) []*testCase {
	for _, test := range cases {
		test.skipped = true
	}
	return cases
}

func (this *fixtureRunner) runFixtureSetup() {
	if this.fixtureSetup >= 0 {
		this.fixture.Method(this.fixtureSetup).Call(nil)
	}
}

func (this *fixtureRunner) runFixtureTeardown() {
	if this.fixtureTeardown >= 0 {
		this.fixture.Method(this.fixtureTeardown).Call(nil)
	}
}

func (this *fixtureRunner) setInnerFixture(innerFixture *Fixture) {
	this.fixture.Elem().FieldByName("Fixture").Set(reflect.ValueOf(innerFixture))
}

func (this *fixtureRunner) leftTestsHandler() {
	if this.skipLeftTests {
		for _, test := range this.tests {
			this.outerT.Run(test.description, test.skip)
		}
		this.outerT.Fail()
	}
}
