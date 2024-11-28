package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func CreateRecorderAndRequest(method string, path string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(method, path, nil)
}

func CreateRecorderAndRequestWithBody(method string, path string, body any) (*httptest.ResponseRecorder, *http.Request) {
	bodyBytes, _ := json.Marshal(body)
	return httptest.NewRecorder(), httptest.NewRequest(method, path, strings.NewReader(string(bodyBytes)))
}

func TestCreateGenericGetManyHandler(t *testing.T) {
	w, r := CreateRecorderAndRequest("GET", "/")
	_ = w
	_ = r
}
