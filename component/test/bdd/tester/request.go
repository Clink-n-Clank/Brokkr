package tester

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

func (a *ActorAPI) IWillUseThisBaseHostName(base string) error {
	baseURL, baseURLErr := url.Parse(base)
	if baseURLErr != nil {
		return baseURLErr
	}

	a.HTTPBaseHost = *baseURL

	return nil
}

func (a *ActorAPI) IWillUseThisBaseHostNameReadEnv(hostEnv string) error {
	baseHost := strings.TrimSpace(os.Getenv(hostEnv))
	if len(baseHost) < 4 { // it can be like m.me, x.com
		return fmt.Errorf("base host name not found or too small")
	}

	baseURL, baseURLErr := url.Parse(baseHost)
	if baseURLErr != nil {
		return baseURLErr
	}

	a.HTTPBaseHost = *baseURL

	return nil
}

func (a *ActorAPI) ISetHeaderWithValue(name string, value string) error {
	a.HTTPHeaders[name] = value

	return nil
}

func (a *ActorAPI) ISendSimpleRequestTo(method, endpoint string) error {
	return a.handleRequest(method, endpoint, nil)
}

func (a *ActorAPI) ISendRequestWithPayload(method, endpoint string, payloadFile string) error {
	payload, payloadErr := loadFile(payloadFile)
	if payloadErr != nil {
		return payloadErr
	}

	return a.handleRequest(method, endpoint, strings.NewReader(string(payload)))
}

func (a *ActorAPI) IAddURLPath(path string) error {
	a.HTTPPath = append(a.HTTPPath, path)

	return nil
}

func (a *ActorAPI) IAddURLPathByStoredKey(key string) error {
	if _, ok := a.Storage[key]; !ok {
		return fmt.Errorf("key %s was not found in memory to use it in query", key)
	}

	a.HTTPPath = append(a.HTTPPath, a.Storage[key])

	return nil
}

func (a *ActorAPI) ISendPreBuiltRequestWithPath(method, endpoint string) error {
	return a.handleRequest(method, endpoint, nil)
}

func (a *ActorAPI) ISendPreBuiltRequestWithPathAndPayload(method, endpoint, payloadFile string) error {
	payload, payloadErr := loadFile(payloadFile)
	if payloadErr != nil {
		return payloadErr
	}

	return a.handleRequest(method, endpoint, strings.NewReader(string(payload)))
}

func (a *ActorAPI) ISendPreBuiltRequestWithStoredPathAndPayload(method, endpoint string) error {
	return a.handleRequest(method, endpoint, strings.NewReader(a.Payload))
}

func (a *ActorAPI) AddQueryRequestParams(dt *godog.Table) error {
	headerRowIndex := 0

	for valueRowIndex := 1; valueRowIndex < len(dt.Rows); valueRowIndex++ {
		for valueIndex, param := range dt.Rows[headerRowIndex].Cells {
			a.HTTPQuery[param.Value] = dt.Rows[valueRowIndex].Cells[valueIndex].Value
		}
	}

	return nil
}

func (a *ActorAPI) AddQueryRequestParamFromStoredKey(param, key string) error {
	if _, ok := a.Storage[key]; !ok {
		return fmt.Errorf("key %s was not found in memory to use in payload", key)
	}
	a.HTTPQuery[param] = a.Storage[key]

	return nil
}

func (a *ActorAPI) AddQueryRequestParamWithValue(param, value string) error {
	a.HTTPQuery[param] = value

	return nil
}

func (a *ActorAPI) ICreatePayloadWithValue(payloadTemplate, key string) error {
	if _, ok := a.Storage[key]; !ok {
		return fmt.Errorf("key %s was not found in memory to use in payload", key)
	}

	a.Payload = strings.Replace(payloadTemplate, "#value#", a.Storage[key], 1)

	return nil
}

func (a *ActorAPI) ICreatePayload(rawPayload string) error {
	if len(rawPayload) < 2 {
		return fmt.Errorf("data payload (%s) is empty or too small", rawPayload)
	}

	a.Payload = rawPayload

	return nil
}

func (a *ActorAPI) IUseThisPayload(payloadTemplate string) error {
	a.Payload = payloadTemplate

	return nil
}
