package handlers_test

import (
	"net/http"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/handlers"
	testcommon "github.com/CelanMatjaz/job_application_tracker_api/pkg/test_common"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestCreateGetManyTagsFromTableHandler(t *testing.T) {
	handler := handlers.CreateGetManyTagsFromTableHandler(func(_ int, _ int) ([]types.Tag, error) {
		return []types.Tag{}, nil
	})

	t.Run("test without id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})

	t.Run("test with id without path param", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, addId(req))
		testcommon.AssertError(t, err, types.InvalidPathParam)
	})

	t.Run("test with id with path param", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodGet, "/", nil)
		err := handler(res, addPathParam(addId(req), "id", "1"))
		testcommon.AssertNotError(t, err, "")
	})
}
