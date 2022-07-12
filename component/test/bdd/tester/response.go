package tester

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"github.com/xeipuuv/gojsonschema"
)

func (a *ActorAPI) TheResponseCodeShouldBe(code int) error {
	if code != a.HTTPLastResp.ResponseObj.StatusCode {
		return fmt.Errorf(
			"expected response code: %d, but it is: %d\n%s", code,
			a.HTTPLastResp.ResponseObj.StatusCode,
			createHTTPResponseDumpFrom(a.HTTPLastResp),
		)
	}

	return nil
}

func (a *ActorAPI) TheResponseShouldBeEmpty() error {
	if len(a.HTTPLastResp.Body) > 0 {
		return fmt.Errorf(
			"expected response to be empty response but its has data...\n%s",
			createHTTPResponseDumpFrom(a.HTTPLastResp),
		)
	}

	return nil
}

func (a *ActorAPI) TheResponseShouldMatchJsonSchema(path string) error {
	schemaContent, schemaContentErr := loadFile(path)
	if schemaContentErr != nil {
		return schemaContentErr
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaContent)
	documentLoader := gojsonschema.NewStringLoader(string(a.HTTPLastResp.Body))

	result, validationErr := gojsonschema.Validate(schemaLoader, documentLoader)
	if validationErr != nil {
		return fmt.Errorf("unable to load json document, err:%w\n%s", validationErr, createHTTPResponseDumpFrom(a.HTTPLastResp))
	}

	if !result.Valid() {
		var schemaErrors []string
		for _, schemaErr := range result.Errors() {
			schemaErrors = append(schemaErrors, schemaErr.String())
		}

		return fmt.Errorf("the response has other Json schema %s\n %v", path, schemaErrors)
	}

	return nil
}

func (a *ActorAPI) TheResponseShouldMatchJSON(body *godog.DocString) error {
	actual := strings.Trim(string(a.HTTPLastResp.Body), "\n")
	expected := body.Content

	match, matchErr := isSameJson(actual, expected)
	if matchErr != nil {
		return fmt.Errorf(
			"unable to match json, err:%w\n%s",
			matchErr,
			createHTTPResponseDumpFrom(a.HTTPLastResp),
		)
	}

	if !match {
		return fmt.Errorf("expected json %s, does not match actual: %s", expected, actual)
	}

	return nil
}

func (a *ActorAPI) TheResponseShouldHaveValueInRegex(expectedValue, regexPattern string) error {
	e := strings.Replace(expectedValue, "%SNP_USERNAME%", os.Getenv("SNP_USERNAME"), 1)
	pattern := fmt.Sprintf("%s%s", regexPattern, e)
	re := regexp.MustCompile(pattern)

	if !re.Match(a.HTTPLastResp.Body) {
		return fmt.Errorf(
			"expected to find value: %s with regex pattern %s, but it not found in response:\n%s",
			expectedValue,
			regexPattern,
			createHTTPResponseDumpFrom(a.HTTPLastResp),
		)
	}

	return nil
}

func (a *ActorAPI) TheResponseShouldMatchRegex(regexPattern string) error {
	re := regexp.MustCompile(regexPattern)
	if !re.Match(a.HTTPLastResp.Body) {
		return fmt.Errorf(
			"expected to find any value with regex pattern %s, but it is not found in response:\n%s",
			regexPattern,
			createHTTPResponseDumpFrom(a.HTTPLastResp),
		)
	}

	return nil
}
