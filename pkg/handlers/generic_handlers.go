package handlers

import (
	"net/http"
	"strconv"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/go-chi/chi/v5"
)

func CreateGenericGetManyHandler[T any](
	get func(int) ([]T, error),
	sendJson func(w http.ResponseWriter, data []T) error,
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.Context().Value(middleware.AccountIdKey).(int)
		if accountId == 0 {
			return types.Unauthenticated
		}

		data, err := get(accountId)
		if err != nil {
			return err
		}

		return sendJson(w, data)
	}
}

func CreateGenericGetManyWithPaginationHandler[T any](
	get func(int, types.PaginationParams) ([]T, error),
	sendJson func(w http.ResponseWriter, data []T) error,
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.Context().Value(middleware.AccountIdKey).(int)
		if accountId == 0 {
			return types.Unauthenticated
		}

		pagination := types.GetPaginationParams(r)
		data, err := get(accountId, pagination)
		if err != nil {
			return err
		}

		return sendJson(w, data)
	}
}

func CreateGenericGetSingleHandler[T any](
	get func(int, int) (T, error),
	sendJson func(w http.ResponseWriter, data T) error,
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.Context().Value(middleware.AccountIdKey).(int)
		if accountId == 0 {
			return types.Unauthenticated
		}

		recordIdParam := chi.URLParam(r, "id")
		recordId, err := strconv.Atoi(recordIdParam)
		if err != nil {
			return types.InvalidPathParam
		}

		data, err := get(accountId, recordId)
		if err != nil {
			return err
		}

		return sendJson(w, data)
	}
}

func CreateGenericPostHandler[T any](
	create func(int, T) (T, error),
	sendJson func(w http.ResponseWriter, data T) error,
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.Context().Value(middleware.AccountIdKey).(int)
		if accountId == 0 {
			return types.Unauthenticated
		}

		body, err := decodeAndValidateBody[T](r)
		if err != nil {
			return err
		}

		data, err := create(accountId,body)
		if err != nil {
			return err
		}

		return sendJson(w, data)
	}
}

func CreateGenericPutHandler[T any](
	update func(int, int, T) (T, error),
	sendJson func(w http.ResponseWriter, data T) error,
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.Context().Value(middleware.AccountIdKey).(int)
		if accountId == 0 {
			return types.Unauthenticated
		}

		recordIdParam := chi.URLParam(r, "id")
		recordId, err := strconv.Atoi(recordIdParam)
		if err != nil {
			return types.InvalidPathParam
		}

		body, err := decodeAndValidateBody[T](r)
		if err != nil {
			return err
		}

		data, err := update(accountId, recordId, body)
		if err != nil {
			return err
		}

		return sendJson(w, data)
	}
}

func CreateGenericDeleteHandler(
	delete func(int, int) error,
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.Context().Value(middleware.AccountIdKey).(int)
		if accountId == 0 {
			return types.Unauthenticated
		}

		recordIdParam := chi.URLParam(r, "id")
		recordId, err := strconv.Atoi(recordIdParam)
		if err != nil {
			return types.InvalidPathParam
		}

		err = delete(accountId, recordId)
		if err != nil {
			return err
		}

		return nil
	}
}
