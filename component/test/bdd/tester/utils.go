package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

func hideSensitiveString(str *string, hide []string) {
	for _, secret := range hide {
		var re = regexp.MustCompile(secret)
		*str = re.ReplaceAllString(*str, "<sensitive>")
	}
}

func isSameJson(j1, j2 string) (bool, error) {
	var o1, o2 interface{}

	if err := json.Unmarshal([]byte(j1), &o1); err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(j2), &o2); err != nil {
		return false, err
	}

	return reflect.DeepEqual(o1, o2), nil
}

func createHTTPResponseDumpFrom(r HTTPResponse) string {
	dump := "-----Original HTTP Dump-----"

	if r.ResponseObj.Request != nil && r.ResponseObj.Request.URL != nil {
		dump = fmt.Sprintf("%s\nRequest URL: %s", dump, r.ResponseObj.Request.URL.String())
		dump = fmt.Sprintf("%s\nRequest Headers: %s", dump, r.ResponseObj.Request.Header)
	}

	dump = fmt.Sprintf("%s\nResponse Headers: %s", dump, r.ResponseObj.Header)
	dump = fmt.Sprintf("%s\nResponse ContentLength: %d", dump, r.ResponseObj.ContentLength)
	dump = fmt.Sprintf("%s\nResponse Status: %s", dump, r.ResponseObj.Status)
	dump = fmt.Sprintf("%s\nResponse Body len: %d\nBody: %v", dump, len(r.Body), string(r.Body))

	return dump
}

func loadFile(path string) ([]byte, error) {
	path = strings.Trim(path, "/")

	absPath, absErr := filepath.Abs(path)
	if absErr != nil {
		return []byte{}, fmt.Errorf("unable to resolve file path: %s", absPath)
	}

	if _, osStatErr := os.Stat(absPath); os.IsNotExist(osStatErr) {
		return []byte{}, fmt.Errorf("file does not exist: %s", absPath)
	}

	content, readFileErr := ioutil.ReadFile(absPath)
	if readFileErr != nil {
		return []byte{}, fmt.Errorf("cannot open file: %w", readFileErr)
	}

	return content, nil
}
