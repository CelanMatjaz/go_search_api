package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	testcommon "github.com/CelanMatjaz/job_application_tracker_api/pkg/test_common"
)

func createRequest(t *testing.T, s *httptest.Server, method string, path string, body any) *http.Request {
	var reader *strings.Reader
	if body == nil {
		reader = &strings.Reader{}
	} else {
		bodyJson, err := json.Marshal(body)
		testcommon.AssertNotError(t, err, "error marshalling json body")
		reader = strings.NewReader(string(bodyJson))
	}

	req := httptest.NewRequest(method, strings.Join([]string{s.URL, path}, "/"), reader)
	return req
}

func getResponse(t *testing.T, req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)
	testcommon.AssertNotError(t, err, "error making request")
	return res
}

func getBodyString(t *testing.T, res *http.Response) string {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading body, %s", err.Error())
	}

	return string(resBody)
}

func getBody[T any](t *testing.T, res *http.Response) T {
	var body T
	err := json.Unmarshal([]byte(getBodyString(t, res)), &body)
	testcommon.AssertNotError(t, err, "error unmarshalling body")

	return body
}

func newRequestAndRecorder(t *testing.T, method string, path string, body any) (*httptest.ResponseRecorder, *http.Request) {
	res := httptest.NewRecorder()

	if body != nil {
		bodyJson, err := json.Marshal(body)
		testcommon.AssertNotError(t, err, "error marshalling json body")
		req := httptest.NewRequest(method, path, bytes.NewReader(bodyJson))
		req.Header.Set("Content-Type", "application/json")
		req.Close = true
		return res, req
	}

	return res, httptest.NewRequest(method, path, nil)
}

func newRequestWithContextValues(req *http.Request, key any, value any) *http.Request {
	ctx := context.WithValue(req.Context(), key, value)
	return req.WithContext(ctx)
}

func assertResponseStatus(t *testing.T, res *httptest.ResponseRecorder, statusCode int) {
	result := res.Result()
	testcommon.Assert(
		t, result.StatusCode == statusCode,
		"response status code does not equal expected status code\nexpected: %d\nactual:   %d",
		statusCode,
		result.StatusCode)
}
