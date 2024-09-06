package db

import (
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)


type AuthStore interface {
	GetInternalUserById(id int) (types.InternalUser, error)
	GetInternalUserByEmail(email string) (types.InternalUser, error)
	CreateUser(user types.InternalUser) (types.InternalUser, error)
}

type ResumeStore1 interface {
	GetResumes(userId int, paginationParams service.PaginationParams) ([]types.Resume, error)
	GetResume(id int) (types.Resume, error)
	CreateResume(types.Resume) (types.Resume, error)
	UpdateResume(types.Resume) (types.Resume, error)
	DeleteResume(types.Resume) error
}

type ResumeTagStore interface {
	GetResumeTags(resumeId int) ([]types.ResumeTag, error)
	CreateResumeTag(types.ResumeTag) (types.ResumeTag, error)
	UpdateResumeTag(types.ResumeTag) (types.ResumeTag, error)
	DeleteResumeTag(types.ResumeTag) error
}
