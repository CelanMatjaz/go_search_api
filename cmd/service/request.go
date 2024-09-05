package service

import (
	"encoding/json"
	"net/http"
	"time"
)

type response struct {
	Timestamp  time.Time `json:"timestamp"`
	StatusCode int       `json:"status_code"`
}

type JsonResponse struct {
	Data any `json:"data"`
	response
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
	response
}

func SendJsonResponse(w http.ResponseWriter, data any, statusCode int) {
	r := JsonResponse{
		Data: data,
		response: response{
			Timestamp:  time.Now(),
			StatusCode: statusCode,
		},
	}

	sendJson(w, r, statusCode)
}

func SendErrorsResponse(w http.ResponseWriter, errors []string, statusCode int) {
	r := ErrorResponse{
		Errors: errors,
		response: response{
			Timestamp:  time.Now(),
			StatusCode: statusCode,
		},
	}

	sendJson(w, r, statusCode)
}

func SendInternalServerError(w http.ResponseWriter) {
	SendErrorsResponse(w, []string{"Internal server error"}, http.StatusInternalServerError)
}

func sendJson(w http.ResponseWriter, val any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(val)
}
