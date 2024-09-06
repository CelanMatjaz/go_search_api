package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var UserIdKey = "USER_ID"

type response struct {
	Timestamp  time.Time         `json:"timestamp"`
	StatusCode int               `json:"status_code"`
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

type PaginationParams struct {
	Page   int `json:"page"`
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

func (p *PaginationParams) GetOffset() int {
	return p.Offset + (p.Page-1)*p.Count
}

func GetPaginationParams(r *http.Request) PaginationParams {
	return PaginationParams{
		Page:   customAtoi(r.URL.Query().Get("page"), 1, 1),
		Offset: customAtoi(r.URL.Query().Get("offset"), 0, 0),
		Count:  customAtoiClamp(r.URL.Query().Get("count"), 10, 1, 100),
	}
}

func customAtoi(str string, defaultValue int, minValue int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}

	val = max(val, minValue)

	return val
}

func customAtoiClamp(str string, defaultValue int, minValue int, maxValue int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}

	val = max(min(val, maxValue), minValue)

	return val
}
