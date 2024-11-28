package handlers

import (
	"net/http"
	"strconv"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
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
		r.Get("/application-presets/{id}", CreateHandler(CreateGetManyTagsFromTableHandler(h.store.GetApplicationPresetTags)))
		r.Get("/application-sections/{id}", CreateHandler(CreateGetManyTagsFromTableHandler(h.store.GetApplicationSectionTags)))
		r.Get("/resume-presets/{id}", CreateHandler(CreateGetManyTagsFromTableHandler(h.store.GetResumePresetTags)))
		r.Get("/resume-sections/{id}", CreateHandler(CreateGetManyTagsFromTableHandler(h.store.GetResumeSectionTags)))

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

func CreateGetManyTagsFromTableHandler(
	get func(int, int) ([]types.Tag, error),
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId, err := getAccountId(r)
		if err != nil {
			return err
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
