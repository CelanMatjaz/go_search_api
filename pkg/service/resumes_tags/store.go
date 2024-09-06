package resumestags

import (
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func NewStore(connection *db.DbConnection) *db.GenericStore[types.ResumeTag] {
	tableName := "resume_tags"
	scanFields := []string{"id", "resume_id", "label"}
	neededFields := []string{"resume_id", "label"}
	updateFields := []string{"label"}

	return &db.GenericStore[types.ResumeTag]{
		Db:              connection.DB,
		Scanner:         &resumeTagScanner{},
		SelectManyQuery: db.CreateSelectManyQuery(tableName, scanFields),
		CreateQuery:     db.CreateCreateQuery(tableName, neededFields, scanFields),
		UpdateQuery:     db.CreateUpdateQuery(tableName, updateFields, scanFields),
		DeleteQuery:     db.CreateDeleteQuery(tableName),
	}
}

type resumeTagScanner struct{}

func (s *resumeTagScanner) Scan(row db.Scannable) (types.ResumeTag, error) {
	var r types.ResumeTag
	return r, row.Scan(
		&r.Id,
		&r.ResumeId,
		&r.Label,
	)
}
