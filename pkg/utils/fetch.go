package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func MakeRequest[ResponseBody any](url string, method string, headers map[string]string, requestBody []byte) (ResponseBody, int, error) {
	var body ResponseBody
	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return body, 0, err
	}

	for key, val := range headers {
		request.Header.Set(key, val)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return body, 0, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return body, 0, err
	}

	err = json.Unmarshal(responseBody, &body)
	return body, response.StatusCode, err
}
