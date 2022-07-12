package integration_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/Clink-n-Clank/Brokkr/component/test/bdd"
)

func TestMain(m *testing.M) {
	bddTester := bdd.NewBDDTester("integration_tests")

	bddTester.AddScenario(bdd.Scenario{
		Expr:     `^This test will pass$`,
		StepFunc: testIsPassing,
	})

	bddTester.AddScenario(bdd.Scenario{
		Expr:     `^This test will fail`,
		StepFunc: testIsFailing,
	})

	status := bddTester.RunTests()
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func testIsPassing() error {
	return nil
}

func testIsFailing() error {
	return fmt.Errorf("test failing")
}
