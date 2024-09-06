package types

type Resume struct {
	Common
	UserId int    `json:"user_id" db:"user_id"`
	Name   string `json:"name" db:"name"`
	Note   string `json:"note" db:"note"`
	Timestamps
}

type ResumeTag struct {
	Common
	ResumeId int    `json:"resume_id" db:"resume_id"`
	Label    string `json:"label" db:"label"`
}
