package types

import (
	"net/http"
	"strconv"
)

type ApiError struct {
	StatusCode int
	Errors     []string
	error      error
}

func (err ApiError) Error() string {
	return err.error.Error()
}

func CreateApiError(err error, errors []string, statusCode int) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Errors:     errors,
		error:      err,
	}
}

var (
	UnparsableJsonBody   = CreateApiError(InvalidJsonBodyErr, UnparseableJsonBodyErrors, http.StatusUnprocessableEntity)
	InvalidJsonBody      = CreateApiError(InvalidJsonBodyErr, InvalidJsonBodyErrors, http.StatusUnprocessableEntity)
	AccountAlreadyExists = CreateApiError(AccountAlreadyExistsErr, AccountAlreadyExistsErrors, http.StatusConflict)
	AccountDoesNotExist  = CreateApiError(AccountDoesNotExistErr, AccountAlreadyExistsErrors, http.StatusUnauthorized)
	InvalidPassword      = CreateApiError(InvalidJsonBodyErr, InvalidPasswordErrors, http.StatusBadRequest)
	Unauthenticated      = CreateApiError(UnauthenticatedErr, UnauthenticatedErrors, http.StatusUnauthorized)
	UnknownOAuthProvider = CreateApiError(UnknownOAuthProviderErr, UnknownOAuthProviderErrors, http.StatusNotFound)
	OAuthProviderIssues  = CreateApiError(OAuthProviderIssuesErr, OAuthProviderIssueErrors, http.StatusServiceUnavailable)
	UnverifiedOAuthEmail = CreateApiError(OAuthUnverifiedEmailErr, UnverifiedOAuthEmailErrors, http.StatusBadRequest)
	InvalidPathParam     = CreateApiError(InvalidPathParamErr, UnauthenticatedErrors, http.StatusBadRequest)
	PasswordsDoNotMatch  = CreateApiError(InvalidJsonBodyErr, PasswordsDoNotMathErrors, http.StatusBadRequest)
)

const (
	defaultPageSize = 10
	minPageSize     = 1
	maxPageSize     = 100
)

type PaginationParams struct {
	Page   int `json:"page"`
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

func DefaultPagaintion() PaginationParams {
	return PaginationParams{
		Page:   1,
		Count:  10,
		Offset: 0,
	}
}

func (p *PaginationParams) GetOffset() int {
	return p.Offset + (p.Page-1)*p.Count
}

func GetPaginationParams(r *http.Request) PaginationParams {
	return PaginationParams{
		Page:   customAtoi(r.URL.Query().Get("page"), 1, 1),
		Offset: customAtoi(r.URL.Query().Get("offset"), 0, 0),
		Count:  customAtoiClamp(r.URL.Query().Get("count"), defaultPageSize, minPageSize, maxPageSize),
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
