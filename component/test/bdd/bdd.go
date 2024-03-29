package bdd

import (
	"github.com/Clink-n-Clank/Brokkr/component/test/bdd/tester"
	"github.com/cucumber/godog"
	"github.com/spf13/pflag"
)

var (
	opts = godog.Options{}

	ActorContext *tester.ActorAPI
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

// Scenario contains an expression and a step function that it triggers
type Scenario struct {
	Expr     interface{}
	StepFunc interface{}
}

// Tester has methods for adding scenarios and running tests
type Tester struct {
	name      string
	scenarios []Scenario
}

// NewBDDTester returns a new bdd tester where you can add test scenarios and run them
func NewBDDTester(name string) *Tester {
	t := &Tester{
		name:      name,
		scenarios: []Scenario{},
	}

	ActorContext = tester.NewActorAPI()
	preloadScenarioActions(t, ActorContext)

	return t
}

// AddScenario adds a new scenario
func (t *Tester) AddScenario(s Scenario) {
	isRelacedDefault := false
	for i, old := range t.scenarios {
		if old.Expr == s.Expr {
			t.scenarios[i] = s
			isRelacedDefault = true
			break
		}
	}

	if !isRelacedDefault {
		t.scenarios = append(t.scenarios, s)
	}
}

// RunTests will run the tests and match them to the given scenarios
func (t *Tester) RunTests() int {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                t.name,
		ScenarioInitializer: t.initializeScenario,
		Options:             &opts,
	}.Run()

	return status
}

func (t *Tester) initializeScenario(ctx *godog.ScenarioContext) {
	for _, scenario := range t.scenarios {
		ctx.Step(scenario.Expr, scenario.StepFunc)
	}
}
