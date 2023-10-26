package advanced

import (
	"fmt"
	"github.com/bugVanisher/gunit"
	"testing"
)

func TestTableDriverFixture(t *testing.T) {
	gunit.Run(new(TableDriverFixture), t)
}

type TableDriverFixture struct {
	*gunit.Fixture
}

func (f *TableDriverFixture) Setup() {
}

func (f *TableDriverFixture) FixtureSetup() {
}

func (f *TableDriverFixture) FixtureTeardown() {
}

func (f *TableDriverFixture) Test() {
	tests := []struct {
		x      int
		y      int
		result int
	}{
		{1, 2, 3},
		{-1, -1, -2},
		{0, 0, 0},
		{100000, 1, 100001},
		{99999, 1, 100000},
		{666, 0, 666},
		{996, 1, 997},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%d+%d", test.x, test.y)
		f.Run("", func(fixture *gunit.Fixture) {
			fixture.GetLogger().Info().Str("name", fixture.T().Name()).Msg(name)
			if result := Add(test.x, test.y); result != test.result {
				fixture.Errorf("expected %d but got %d", test.result, result)
			}
		})
	}
}

func Add(x, y int) int {
	return x + y
}
