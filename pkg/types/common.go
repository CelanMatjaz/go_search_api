package types

import "time"

type Common struct {
	Id int `json:"id" db:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func (c Common) GetId() int {
	return c.Id
}
