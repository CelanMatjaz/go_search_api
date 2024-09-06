package resumestags

import (
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
)

type ResumeTagPostBody struct {
	Label *string `json:"label"`
}

func (l *ResumeTagPostBody) IsValid() error {
	if l.Label == nil {
		return service.InvalidBodyErr
	}

	return nil
}
