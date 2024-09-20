package db

import "github.com/CelanMatjaz/job_application_tracker_api/pkg/types"

type Store interface {
	AuthStore
	TagStore
	ApplicationStore
}

type AuthStore interface {
	GetAccountById(id int) (*types.Account, error)
	GetAccountByEmail(email string) (*types.Account, error)
	CreateAccount(account types.Account) (*types.Account, error)
	CreateAccountWithOAuth(account types.Account, tokenResponse types.TokenResponse, clientId int) (*types.Account, error)
	UpdateAccountToOAuth(account types.Account, tokenResponse types.TokenResponse, clientId int) error
	GetOAuthClientByName(name string) (*types.OAuthClient, error)
}

type TagStore interface {
	GetTags(accountId int) ([]types.Tag, error)
	GetTag(accountId int, tagId int) (*types.Tag, error)
	CreateTag(accountId int, tag types.TagBody) (*types.Tag, error)
	UpdateTag(accountId int, tagId int, tag types.TagBody) (*types.Tag, error)
	DeleteTag(accountId int, tagId int) error

	GetApplicationPresetTags(accountId int, applicationId int) ([]types.Tag, error)
	GetApplicationSectionTags(accountId int, applicationSectionId int) ([]types.Tag, error)
	GetResumePresetTags(accountId int, resumeId int) ([]types.Tag, error)
	GetResumeSectionTags(accountId int, resumeSectionId int) ([]types.Tag, error)
}

type ApplicationStore interface {
	GetApplicationPresets(accountId int) ([]types.ApplicationPreset, error)
	GetApplicationPreset(accountId int, presetId int) (*types.ApplicationPreset, error)
	CreateApplicationPreset(accountId int, presetId types.ApplicationPresetBody) (*types.ApplicationPreset, error)
	UpdateApplicationPreset(accountId int, presetId int, preset types.ApplicationPresetBody) (*types.ApplicationPreset, error)
	DeleteApplicationPreset(accountId int, presetId int) error

	GetApplicationSections(accountId int) ([]types.ApplicationSection, error)
	GetApplicationSection(accountId int, presetId int) (*types.ApplicationSection, error)
	CreateApplicationSection(accountId int, presetId types.ApplicationSectionBody) (*types.ApplicationSection, error)
	UpdateApplicationSection(accountId int, presetId int, preset types.ApplicationSectionBody) (*types.ApplicationSection, error)
	DeleteApplicationSection(accountId int, presetId int) error
}

type ResumeStore interface {
	GetResumePresets(accountId int) ([]types.ResumePreset, error)
	GetResumePreset(accountId int, presetId int) (*types.ResumePreset, error)
	CreateResumePreset(accountId int, presetId types.ResumePresetBody) (*types.ResumePreset, error)
	UpdateResumePreset(accountId int, presetId int, preset types.ResumePresetBody) (*types.ResumePreset, error)
	DeleteResumePreset(accountId int, presetId int) error

	GetResumeSections(accountId int) ([]types.ResumeSection, error)
	GetResumeSection(accountId int, presetId int) (*types.ResumeSection, error)
	CreateResumeSection(accountId int, presetId types.ResumeSectionBody) (*types.ResumeSection, error)
	UpdateResumeSection(accountId int, presetId int, preset types.ResumeSectionBody) (*types.ResumeSection, error)
	DeleteResumeSection(accountId int, presetId int) error
}
