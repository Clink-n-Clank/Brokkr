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

type Scenario struct {
	expr interface{}
	stepFunc interface{}
}

type bddTester struct {
	name string
	scenarios []Scenario
}

func NewBDDTester(name string) *bddTester {
	return &bddTester{
		name:            name,
		scenarios: []Scenario{},
	}
}

func (t *bddTester) AddScenario(s Scenario) {
	t.scenarios = append(t.scenarios, s)
}

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
		ctx.Step(scenario.expr, scenario.stepFunc)
	}
}
