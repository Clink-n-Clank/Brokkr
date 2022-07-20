package bdd

import "github.com/Clink-n-Clank/Brokkr/component/test/bdd/tester"

func preloadScenarioActions(t *Tester, a *tester.ActorAPI) {
	//
	// Make a request
	//
	t.AddScenario(Scenario{
		Expr:     `^I will use this base host name "(.*)"$`,
		StepFunc: a.IWillUseThisBaseHostName,
	})
	t.AddScenario(Scenario{
		Expr:     `^I will use this base host name from env var "(.*)"$`,
		StepFunc: a.IWillUseThisBaseHostNameReadEnv,
	})
	t.AddScenario(Scenario{
		Expr:     `^I set header "([^"]*)" with value "([^"]*)"$`,
		StepFunc: a.ISetHeaderWithValue,
	})
	t.AddScenario(Scenario{
		Expr:     `^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`,
		StepFunc: a.ISendSimpleRequestTo,
	})
	t.AddScenario(Scenario{
		Expr:     `^I send "([^"]*)" request to "([^"]*)" with payload "([^"]*)"$`,
		StepFunc: a.ISendRequestWithPayload,
	})
	//
	// Request constructor
	//
	t.AddScenario(Scenario{
		Expr:     `^I will add URL path "([^"]*)"$`,
		StepFunc: a.IAddURLPath,
	})
	t.AddScenario(Scenario{
		Expr:     `^I will add URL path value by stored key "(#[^"]*#)"$`,
		StepFunc: a.IAddURLPathByStoredKey,
	})
	t.AddScenario(Scenario{
		Expr:     `^Now I will send "(GET|POST|PUT|DELETE|PATCH)" request to URL "([^"]*)" with pre-built path$`,
		StepFunc: a.ISendPreBuiltRequestWithPath,
	})
	t.AddScenario(Scenario{
		Expr:     `^Now I will send "(GET|POST|PUT|DELETE|PATCH)" request to URL "([^"]*)" with pre-built path and with payload "([^"]*)"$`,
		StepFunc: a.ISendPreBuiltRequestWithPathAndPayload,
	})
	t.AddScenario(Scenario{
		Expr:     `^Now I will send "(GET|POST|PUT|DELETE|PATCH)" request to URL "([^"]*)" with pre-built path and pre-built payload$`,
		StepFunc: a.ISendPreBuiltRequestWithStoredPathAndPayload,
	})
	// Query
	t.AddScenario(Scenario{
		Expr:     `^Use the following query params to the request$`,
		StepFunc: a.AddQueryRequestParams,
	})
	t.AddScenario(Scenario{
		Expr:     `^I will create query param "(.*)" with value from stored key "(#[^"]*#)"$`,
		StepFunc: a.AddQueryRequestParamFromStoredKey,
	})
	t.AddScenario(Scenario{
		Expr:     `^I will create query param "(.*)" with value "(.*)"$`,
		StepFunc: a.AddQueryRequestParamWithValue,
	})
	//
	// Payload / Request Body
	//
	t.AddScenario(Scenario{
		Expr:     `^I will create payload "({.*#value#.*})" with value from stored key "(#[^"]*#)"$`,
		StepFunc: a.ICreatePayloadWithValue,
	})
	t.AddScenario(Scenario{
		Expr:     `^I will create payload: "(.*)"$`,
		StepFunc: a.ICreatePayload,
	})
	t.AddScenario(Scenario{
		Expr:     `^I will use this payload: "(.*)"$`,
		StepFunc: a.IUseThisPayload,
	})
	//
	// Response validation
	//
	t.AddScenario(Scenario{
		Expr:     `^the response code should be (\d+)$`,
		StepFunc: a.TheResponseCodeShouldBe,
	})
	t.AddScenario(Scenario{
		Expr:     `^the response should be empty$`,
		StepFunc: a.TheResponseShouldBeEmpty,
	})
	t.AddScenario(Scenario{
		Expr:     `^the response should match json schema "([^"]*)"$`,
		StepFunc: a.TheResponseShouldMatchJsonSchema,
	})
	t.AddScenario(Scenario{
		Expr:     `^the response should match json:$`,
		StepFunc: a.TheResponseShouldMatchJSON,
	})
	t.AddScenario(Scenario{
		Expr:     `^the response should have "([^"]*)" in regex:(.*)$`,
		StepFunc: a.TheResponseShouldHaveValueInRegex,
	})
	t.AddScenario(Scenario{
		Expr:     `^the response should match regex:(.*)$`,
		StepFunc: a.TheResponseShouldMatchRegex,
	})
	//
	// Miscellaneous
	//
	t.AddScenario(Scenario{
		Expr:     `^Current time is after "(.*)" I will do`,
		StepFunc: a.WhenCurrentTimeAfter,
	})
	t.AddScenario(Scenario{
		Expr:     `^Current time is before "(.*)" I will do`,
		StepFunc: a.WhenCurrentTimeBefore,
	})
	t.AddScenario(Scenario{
		Expr:     `^I will wait "([^"]*)" seconds$`,
		StepFunc: a.IWaitSeconds,
	})
	t.AddScenario(Scenario{
		Expr:     `^Set concurrent requests count "(\d+)" with delay "(\d+)" milliseconds$`,
		StepFunc: a.SetConcurrentRequestsAmount,
	})
	t.AddScenario(Scenario{
		Expr:     `^Print Requests Metrics$`,
		StepFunc: a.PrintMetricsResult,
	})
	t.AddScenario(Scenario{
		Expr:     `^Percentage of dropped request must be less than "(.*)"$`,
		StepFunc: a.MetricsRequestDropIsLowerThan,
	})
	// Debug
	t.AddScenario(Scenario{
		Expr:     `^Debug - Show response$`,
		StepFunc: a.PrintLastResponse,
	})
}
