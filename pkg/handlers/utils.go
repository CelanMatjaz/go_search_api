package handlers

import (
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

func decodeAndValidateBody[T any](r *http.Request) (T, error) {
	body, err := utils.DecodeJsonBody[T](r)
	if err != nil {
		return body, err
	}

	if validateErrors := utils.Validate(body); len(validateErrors) > 0 {
		return body, types.CreateApiError(types.InvalidJsonBodyErr, validateErrors, http.StatusBadRequest)
	}

	return body, nil
}
