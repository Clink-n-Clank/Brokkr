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
	// Replace action
	bddTester.AddScenario(bdd.Scenario{
		Expr:     `^I set header "([^"]*)" with value "([^"]*)"$`,
		StepFunc: testReplacedHttpQueryAction,
	})

	status := bddTester.RunTests()
	if st := m.Run(); st > status {
		status = st
	}

	v, found := bdd.ActorContext.HTTPQuery["Unit"]
	if !found || v != "Test" {
		fmt.Printf("Expected that scenario action will be replaced (testReplacedHttpQueryAction)\n")
		os.Exit(1)
	}

	os.Exit(status)
}

func testReplacedHttpQueryAction(k, v string) error {
	bdd.ActorContext.HTTPQuery[k] = v

	return nil
}

func testIsPassing() error {
	return nil
}

func testIsFailing() error {
	return fmt.Errorf("test failing")
}
