package types

import "time"

type Common struct {
	Id int `json:"id" db:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TODO: Move 
type JobListing struct {
	Common
	Url     string `json:"url" db:"url"`
	Company string `json:"company" db:"company"`
}

