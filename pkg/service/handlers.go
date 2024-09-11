package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type GenericHttpHandler[Type Validatable] interface {
	GetMultiple(int, PaginationParams) ([]Type, error)
	GetSingle(int, int) (Type, error)
	Create(int, Type) (Type, error)
	Update(int, Type) (Type, error)
	Delete(int, int) (int, error)
}

type Validatable interface {
	IsValid() []string
}

type SendJsonWithPagination[T any] func(w http.ResponseWriter, data []T, p PaginationParams, statusCode int)
type SendJson[T any] func(w http.ResponseWriter, data T, statusCode int)

func CreateGetAllHandler[T Validatable](h GenericHttpHandler[T], sendJson SendJsonWithPagination[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := GetPaginationParams(r)
		userId := r.Context().Value(utils.UserIdKey).(int)

		data, err := h.GetMultiple(userId, pagination)
		if err != nil {
			SendInternalServerError(w)
			return
		}

		sendJson(w, data, pagination, http.StatusOK)
	}
}

var (
	paramNotParsable = []string{"Provided url param was not parsable"}
	notFound         = []string{"Requested item does not exist"}
	notFoundDeleted  = []string{"Requested record to be deleted was not found"}
)

func CreateGetSingleHandler[T Validatable](h GenericHttpHandler[T], sendJson SendJson[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(utils.UserIdKey).(int)
		itemId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			SendErrorsResponse(w, paramNotParsable, http.StatusBadRequest)
			return
		}

		data, err := h.GetSingle(userId, itemId)
		switch err {
		case nil:
			sendJson(w, data, http.StatusOK)
			return
		case types.RecordDoesNotExist:
			SendErrorsResponse(w, notFound, http.StatusNotFound)
			return
		default:
			SendInternalServerError(w)
			return
		}
	}
}

func CreatePostHandler[T Validatable](h GenericHttpHandler[T], sendJson SendJson[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var body T
		decoder.Decode(&body)

		if errors := body.IsValid(); len(errors) > 0 {
			SendErrorsResponse(w, errors, http.StatusBadRequest)
			return
		}

		userId := r.Context().Value(utils.UserIdKey).(int)
		data, err := h.Create(userId, body)

		switch err {
		case nil:
			sendJson(w, data, http.StatusOK)
			return
		default:
			SendInternalServerError(w)
			return
		}
	}
}

func CreatePutHandler[T Validatable](h GenericHttpHandler[T], sendJson SendJson[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var body T
		decoder.Decode(&body)

		if errors := body.IsValid(); len(errors) > 0 {
			SendErrorsResponse(w, errors, http.StatusBadRequest)
			return
		}

		userId := r.Context().Value(utils.UserIdKey).(int)
		data, err := h.Update(userId, body)
		switch err {
		case nil:
			sendJson(w, data, http.StatusOK)
			return
		default:
			SendInternalServerError(w)
			return
		}
	}
}

func CreateDeleteHandler[T Validatable](h GenericHttpHandler[T], sendJson SendJson[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			SendErrorsResponse(w, paramNotParsable, http.StatusBadRequest)
			return
		}

		userId := r.Context().Value(utils.UserIdKey).(int)
		deletedCount, err := h.Delete(userId, itemId)
		if deletedCount == 0 {
			SendErrorsResponse(w, notFoundDeleted, http.StatusNotFound)
			return
		}

		if err != nil {
			SendInternalServerError(w)
			return
		}

		SendJsonResponse(w, nil, http.StatusOK)
	}
}
