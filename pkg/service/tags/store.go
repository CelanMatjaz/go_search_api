package tags

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(connection *db.DbConnection) *Store {
	return &Store{db: connection.DB}
}

func (s *Store) GetUserTag(userId int, tagData int) (types.Tag, error)
func (s *Store) CreateUserTag(userId int, tagData types.Tag) (types.Tag, error)
func (s *Store) UpdateUserTag(userId int, tagData types.Tag) (types.Tag, error)
func (s *Store) DeleteUserTag(userId int, tagId int) error

func (s *Store) GetApplicationTags(userId int, applicationId int) (types.Tag, error)
func (s *Store) GetApplicationSectionTags(userId int, sectionId int) (types.Tag, error)
func (s *Store) GetResumeTags(userId int, resume int) (types.Tag, error)
func (s *Store) GetResumeSectionTags(userId int, sectionId int) (types.Tag, error)
