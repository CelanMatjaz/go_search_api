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

// Since this is used for testing, no need to defer res.Body.Close()
func assertResponseStatus(t *testing.T, req *http.Request, expectedStatus int) *http.Response {
	res := getResponse(t, req)
	testcommon.Assert(
		t, res.StatusCode == expectedStatus,
		"response status does not match expected status\nexpected: %d\nactual:   %d (%s)", expectedStatus, res.StatusCode, res.Status,
	)
	return res
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

// func assertStatusWithActualHandler(t *testing.T, req *http.Request, expectedStatusCode int, handler handlers.HandlerFunc) {
// 	res := httptest.NewRecorder()
// 	handlers.CreateHandler(handler)(res, req)
//
// 	result := res.Result()
// 	if result.StatusCode != expectedStatusCode {
// 		body, _ := io.ReadAll(result.Body)
// 		t.Log("response body:", string(body))
// 		t.Fatalf("response status does not match expected status\nexpected: %d\nactual:   %d (%s)", expectedStatusCode, result.StatusCode, result.Status)
// 	}
// }
