package types

import "time"

type Common struct {
	Id int `json:"id" db:"id"`
}

type WithTags struct {
	Tags []*Tag `json:"tags,omitempty"`
}

type Timestamps struct {
	CreatedAt time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt" json:"updatedAt"`
}

type Verifiable interface {
	Verify() []string
}

func (c Common) GetId() int {
	return c.Id
}

func (wt* WithTags) AppendTag(newTag *Tag) {
	if wt.Tags == nil {
		wt.Tags = make([]*Tag, 0)
	}
	wt.Tags = append(wt.Tags, newTag)
}
