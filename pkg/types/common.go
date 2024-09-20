package types

import "time"

type Common struct {
	Id int `json:"id" db:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt" json:"updatedAt"`
}

type Verifiable interface {
	Verify() []string
}
