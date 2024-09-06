package resumes

import "github.com/CelanMatjaz/job_application_tracker_api/pkg/types"

type ResumePostBody struct {
	Name *string `json:"name"`
	Note *string `json:"note"`
}

func (l *ResumePostBody) IsValid() error {
	if l.Name == nil || l.Note == nil {
		return types.InvalidBodyErr
	}

	return nil
}
