package handlers

import (
	"log"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func CreateHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if error, ok := err.(types.ApiError); ok {
				utils.SendErrors(w, error.Errors, error.StatusCode)
			} else {
				log.Println("Internal server error: ", error.Error())
				utils.SendInternalServerError(w)
			}
            return
		}
	}
}

func decodeAndVerifyBody[T types.Verifiable](r *http.Request) (T, error) {
	body, err := utils.DecodeJsonBody[T](r)
	if err != nil {
		return body, err
	}

	if err := utils.VerifyBody(body); err != nil {
		return body, err
	}

	return body, nil
}
