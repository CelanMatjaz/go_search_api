package db

import (
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type AuthStore interface {
	GetUserById(userId int) (types.User, error)
	GetUserByEmail(email string) (types.User, error)
	CreateUser(user types.User) (types.User, error)
}

type TagStore interface {
	GetUserTag(userId int, tagData int) (types.Tag, error)
	CreateUserTag(userId int, tagData types.Tag) (types.Tag, error)
	UpdateUserTag(userId int, tagData types.Tag) (types.Tag, error)
	DeleteUserTag(userId int, tagId int) (int,error)

	GetApplicationTags(userId int, applicationId int) (types.Tag, error)
	GetApplicationSectionTags(userId int, sectionId int) (types.Tag, error)
	GetResumeTags(userId int, resume int) (types.Tag, error)
	GetResumeSectionTags(userId int, sectionId int) (types.Tag, error)
}

type ApplicationStore interface {
	GetUserApplications(userId int, pagination service.PaginationParams) ([]types.Application, error)
	GetUserApplication(userId int, applicationId int) (types.Application, error)
	CreateUserApplication(userId int, applicationData types.Application) (types.Application, error)
	UpdateUserApplication(userId int, applicationData types.Application) (types.Application, error)
	DeleteUserApplication(userId int, applicationId int) (int,error)

	GetApplicationSections(userId int, pagination service.PaginationParams) ([]types.ApplicationSection, error)
	CreateApplicationSection(userId int, applicationData types.ApplicationSection) (types.ApplicationSection, error)
	UpdateApplicationSection(userId int, applicationData types.ApplicationSection) (types.ApplicationSection, error)
	DeleteApplicationSection(userId int, sectionId int) (int,error)
}

type ResumeStore interface {
	GetUserResumes(userId int, pagination service.PaginationParams) ([]types.Resume, error)
	GetUserResume(userId int, applicationId int) (types.Resume, error)
	CreateUserResume(userId int, applicationData types.Resume) (types.Resume, error)
	UpdateUserResume(userId int, applicationData types.Resume) (types.Resume, error)
	DeleteUserResume(userId int, applicationId int) (int,error)

	GetResumeSections(userId int, pagination service.PaginationParams) ([]types.ResumeSection, error)
	CreateResumesSections(userId int, applicationData types.Resume) (types.Resume, error)
	UpdateResumesSections(userId int, applicationData types.Resume) (types.Resume, error)
	DeleteResumesSections(userId int, sectionId int) (int,error)
}
