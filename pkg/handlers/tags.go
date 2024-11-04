package handlers

import (
	"net/http"
	"strconv"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type TagHandler struct {
	store db.TagStore
}

func NewTagHandler(store db.TagStore) *TagHandler {
	return &TagHandler{store: store}
}

func (h *TagHandler) AddRoutes(r chi.Router) {
	r.Route("/tags", func(r chi.Router) {
		r.Get("/application-presets/{id}", CreateHandler(createGetManyTagsFromTableHandler(h.store.GetApplicationPresetTags)))
		r.Get("/application-sections/{id}", CreateHandler(createGetManyTagsFromTableHandler(h.store.GetApplicationSectionTags)))
		r.Get("/resume-presets/{id}", CreateHandler(createGetManyTagsFromTableHandler(h.store.GetResumePresetTags)))
		r.Get("/resume-sections/{id}", CreateHandler(createGetManyTagsFromTableHandler(h.store.GetResumeSectionTags)))

		r.Get("/", CreateHandler(CreateGenericGetManyHandler(h.store.GetTags, sendTags)))
		r.Get("/{id}", CreateHandler(CreateGenericGetSingleHandler(h.store.GetTag, sendTag)))
		r.Post("/", CreateHandler(CreateGenericPostHandler(h.store.CreateTag, sendTag)))
		r.Put("/{id}", CreateHandler(CreateGenericPutHandler(h.store.UpdateTag, sendTag)))
		r.Delete("/{id}", CreateHandler(CreateGenericDeleteHandler(h.store.DeleteTag)))
	})
}

func sendTags(w http.ResponseWriter, data []types.Tag) error {
	return utils.SendJson(w, struct {
		Tags []types.Tag `json:"tags"`
	}{Tags: data}, http.StatusOK)
}

func sendTag(w http.ResponseWriter, data types.Tag) error {
	return utils.SendJson(w, struct {
		Tag types.Tag `json:"tag"`
	}{Tag: data}, http.StatusOK)
}

func createGetManyTagsFromTableHandler(
	get func(int, int) ([]types.Tag, error),
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

		return sendTags(w, data)
	}
}
