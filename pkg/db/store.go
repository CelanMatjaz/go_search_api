package db

import "github.com/CelanMatjaz/job_application_tracker_api/pkg/types"

type Store interface {
	AuthStore
	TagStore
	ApplicationStore
	ResumeStore
}

type AuthStore interface {
	GetAccountById(id int) (types.Account, bool, error)
	GetAccountByEmail(email string) (types.Account, bool, error)
	CreateAccount(account types.Account) (types.Account, error)
	CreateAccountWithOAuth(account types.Account, tokenResponse types.TokenResponse, clientId int) (types.Account, error)
	UpdateAccountToOAuth(account types.Account, tokenResponse types.TokenResponse, clientId int) error
	GetOAuthClientByName(name string) (types.OAuthClient, bool, error)
}

type TagStore interface {
	GetTags(accountId int) ([]types.Tag, error)
	GetTag(accountId int, tagId int) (types.Tag, error)
	CreateTag(accountId int, tag types.Tag) (types.Tag, error)
	UpdateTag(accountId int, id int, tag types.Tag) (types.Tag, error)
	DeleteTag(accountId int, tagId int) error

	GetApplicationPresetTags(accountId int, applicationId int) ([]types.Tag, error)
	GetApplicationSectionTags(accountId int, applicationSectionId int) ([]types.Tag, error)
	GetResumePresetTags(accountId int, resumeId int) ([]types.Tag, error)
	GetResumeSectionTags(accountId int, resumeSectionId int) ([]types.Tag, error)
}

type ApplicationStore interface {
	GetApplicationPresets(accountId int, pagination types.PaginationParams) ([]types.RecordWithTags[types.ApplicationPreset], error)
	GetApplicationPreset(accountId int, id int) (types.ApplicationPreset, error)
	CreateApplicationPreset(accountId int, body types.ApplicationPreset) (types.ApplicationPreset, error)
	UpdateApplicationPreset(accountId int, id int, body types.ApplicationPreset) (types.ApplicationPreset, error)
	DeleteApplicationPreset(accountId int, id int) error

	GetApplicationSections(accountId int, pagination types.PaginationParams) ([]types.RecordWithTags[types.ApplicationSection], error)
	GetApplicationSection(accountId int, id int) (types.ApplicationSection, error)
	CreateApplicationSection(accountId int, body types.ApplicationSection) (types.ApplicationSection, error)
	UpdateApplicationSection(accountId int, id int, body types.ApplicationSection) (types.ApplicationSection, error)
	DeleteApplicationSection(accountId int, id int) error
}

type ResumeStore interface {
	GetResumePresets(accountId int, pagination types.PaginationParams) ([]types.RecordWithTags[types.ResumePreset], error)
	GetResumePreset(accountId int, id int) (types.ResumePreset, error)
	CreateResumePreset(accountId int, body types.ResumePreset) (types.ResumePreset, error)
	UpdateResumePreset(accountId int, id int, body types.ResumePreset) (types.ResumePreset, error)
	DeleteResumePreset(accountId int, id int) error

	GetResumeSections(accountId int, pagination types.PaginationParams) ([]types.RecordWithTags[types.ResumeSection], error)
	GetResumeSection(accountId int, id int) (types.ResumeSection, error)
	CreateResumeSection(accountId int, body types.ResumeSection) (types.ResumeSection, error)
	UpdateResumeSection(accountId int, id int, body types.ResumeSection) (types.ResumeSection, error)
	DeleteResumeSection(accountId int, id int) error
}
