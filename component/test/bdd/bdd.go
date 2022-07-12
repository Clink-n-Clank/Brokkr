package bdd

import (
	"github.com/Clink-n-Clank/Brokkr/component/test/bdd/tester"
	"github.com/cucumber/godog"
	"github.com/spf13/pflag"
)

var (
	opts         = godog.Options{}
	actorContext *tester.ActorAPI
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

// Scenario contains an expression and a step function that it triggers
type Scenario struct {
	Expr interface{}
	StepFunc interface{}
}

type bddTester struct {
	name string
	scenarios []Scenario
}

// NewBDDTester returns a new bdd tester where you can add test scenarios and run them
func NewBDDTester(name string) *bddTester {
	return &bddTester{
		name:            name,
		scenarios: []Scenario{},
	}
}

// AddScenario adds a new scenario
func (t *bddTester) AddScenario(s Scenario) {
	t.scenarios = append(t.scenarios, s)
}

// RunTests will run the tests and match them to the given scenarios
func (t *bddTester) RunTests() int {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                t.name,
		ScenarioInitializer: t.initializeScenario,
		Options:             &opts,
	}.Run()

	return status
}

func (t *bddTester) initializeScenario(ctx *godog.ScenarioContext) {
	for _, scenario := range t.scenarios {
		ctx.Step(scenario.Expr, scenario.StepFunc)
	}
}