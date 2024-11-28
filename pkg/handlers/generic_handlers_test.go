package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/handlers"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	testcommon "github.com/CelanMatjaz/job_application_tracker_api/pkg/test_common"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func dummySendJson(_ http.ResponseWriter, _ []int) error {
	return nil
}

func addId(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.AccountIdKey, 1)
	return req.WithContext(ctx)
}

func addContextValue(req *http.Request, key any, value any) *http.Request {
	ctx := context.WithValue(req.Context(), key, value)
	return req.WithContext(ctx)
}

func TestCreateGenericGetManyHandler(t *testing.T) {
	handler := handlers.CreateGenericGetManyHandler(func(_ int) ([]int, error) {
		return []int{1, 2, 3, 4, 5}, nil
	}, dummySendJson)

	t.Run("test without id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})

	t.Run("test with id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, addId(req))
		testcommon.AssertNotError(t, err, "")
	})
}

func TestCreateGenericGetManyWithPaginationHandler(t *testing.T) {
	handler := handlers.CreateGenericGetManyWithPaginationHandler(func(_ int, _ types.PaginationParams) ([]int, error) {
		return []int{1, 2, 3, 4, 5}, nil
	}, dummySendJson)

	t.Run("test without id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})

	t.Run("test with id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, addId(req))
		testcommon.AssertNotError(t, err, "")
	})
}

func TestCreateGenericGetSingleHandler(t *testing.T) {
	handler := handlers.CreateGenericGetSingleHandler(func(_ int, _ int) (int, error) {
		return 1, nil
	}, func(_ http.ResponseWriter, _ int) error { return nil })

	t.Run("test without id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})

	t.Run("test with id without param", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, addId(req))
		testcommon.AssertError(t, err, types.InvalidPathParam)
	})

	t.Run("test with id with param", func(t *testing.T) {
		// res, req := newRequestAndRecorder(t, http.MethodPost, "/", nil)
		// err := handler(res, addId(addContextValue(req, "id", 1)))
		// testcommon.AssertNotError(t, err, "")
	})
}

func TestCreateGenericPostHandler(t *testing.T) {
	type GenericBody struct {
		Value string `validate:"required"`
	}

	handler := handlers.CreateGenericPostHandler(func(_ int, _ GenericBody) (GenericBody, error) {
		return GenericBody{}, nil
	}, func(_ http.ResponseWriter, _ GenericBody) error { return nil })

	t.Run("test without id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, "/", nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})

	t.Run("test with id without body", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, "/", nil)
		err := handler(res, addId(req))
		testcommon.AssertError(t, err, types.InvalidJsonBody)
	})

	t.Run("test with id with invalid body", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, "/", GenericBody{})
		err := handler(res, addId(req))
		testcommon.AssertError(t, err, types.InvalidJsonBody)
	})

	t.Run("test with id with valid body", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, "/", GenericBody{Value: "test"})
		err := handler(res, addId(req))
		testcommon.AssertNotError(t, err, "")
	})
}

func TestCreateGenericPutHandler(t *testing.T) {
	type GenericBody struct {
		Value string `validate:"required"`
	}

	handler := handlers.CreateGenericPutHandler(func(_ int, _ int, _ GenericBody) (GenericBody, error) {
		return GenericBody{}, nil
	}, func(_ http.ResponseWriter, _ GenericBody) error { return nil })

	t.Run("test without id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, "/", nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})

	t.Run("test with id without body", func(t *testing.T) {
		// res, req := newRequestAndRecorder(t, http.MethodPost, "/", nil)
		// err := handler(res, addId(req))
		// testcommon.AssertError(t, err, types.InvalidJsonBody)
	})

	t.Run("test with id with invalid body", func(t *testing.T) {
		// res, req := newRequestAndRecorder(t, http.MethodPost, "/", GenericBody{})
		// err := handler(res, addId(req))
		// testcommon.AssertError(t, err, types.InvalidJsonBody)
	})

	t.Run("test with id with valid body", func(t *testing.T) {
		// res, req := newRequestAndRecorder(t, http.MethodPost, "/", GenericBody{Value: "test"})
		// err := handler(res, addId(req))
		// testcommon.AssertNotError(t, err, "")
	})
}

func TestCreateGenericDeleteHandler(t *testing.T) {
	handler := handlers.CreateGenericDeleteHandler(func(_ int, _ int) error {
		return nil
	})

	t.Run("test without id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodDelete, "/", nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})

	t.Run("test with id without path param", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodDelete, "/", nil)
		err := handler(res, addId(req))
		testcommon.AssertError(t, err, types.InvalidPathParam)
	})

	t.Run("test with id with path param", func(t *testing.T) {
		// res, req := newRequestAndRecorder(t, http.MethodDelete, "/", nil)
		// err := handler(res, addId(req))
		// testcommon.AssertNotError(t, err, "")
	})
}
