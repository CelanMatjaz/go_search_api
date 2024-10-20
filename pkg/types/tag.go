package types

type Tag struct {
	Id        NullInt64  `json:"id" db:"id"`
	AccountId NullInt64  `json:"-" db:"account_id" body:""`
	Label     NullString `json:"label" db:"label" body:"" validate:"required,min:1,max:32"`
	Color     NullString `json:"color" db:"color" body:"" validate:"required,len:7"`
}

type TagBody struct {
	Label string `json:"label" db:"label" body:"" validate:"required,min:1,max:32"`
	Color string `json:"color" db:"color" body:"" validate:"required,len:7"`
}

type WithTags struct {
	Tags   []Tag `json:"tags,omitempty"`
	TagIds []int `json:"tagIds,omitempty"`
}

func (wt WithTags) AppendTag(newTag Tag) {
	if wt.Tags == nil {
		wt.Tags = make([]Tag, 0)
	}
	wt.Tags = append(wt.Tags, newTag)
}

func (wti WithTags) GetTagCount() int {
	return len(wti.TagIds)
}
