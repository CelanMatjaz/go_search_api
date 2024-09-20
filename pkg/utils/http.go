package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type response struct {
	Timestamp  time.Time `json:"timestamp"`
	StatusCode int       `json:"statusCode"`
}

type jsonResponse struct {
	Data any `json:"data,omitempty"`
	response
}

type errorsResponse struct {
	Errors []string `json:"errors"`
	response
}

func SendJson(w http.ResponseWriter, data any, statusCode int) error {
	return sendJson(w, jsonResponse{
		Data: data,
		response: response{
			StatusCode: statusCode,
			Timestamp:  time.Now(),
		},
	}, statusCode)
}

func SendErrors(w http.ResponseWriter, errors []string, statusCode int) error {
	return sendJson(w, errorsResponse{
		Errors: errors,
		response: response{
			StatusCode: statusCode,
			Timestamp:  time.Now(),
		},
	}, statusCode)
}

func SendInternalServerError(w http.ResponseWriter) error {
	return SendErrors(w, types.InternalServerErrors, http.StatusInternalServerError)
}

func sendJson(w http.ResponseWriter, val any, statusCode int) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func DecodeJsonBody[T any](r *http.Request) (T, error) {
	decoder := json.NewDecoder(r.Body)
	var body T
	err := decoder.Decode(&body)
	if err != nil {
		return body, types.UnparsableJsonBody
	}

	return body, nil
}

func VerifyBody[T types.Verifiable](body T) error {
	if errors := body.Verify(); len(errors) > 0 {
		return types.CreateApiError(types.InvalidJsonBody, errors, http.StatusUnprocessableEntity)
	}

	return nil
}
