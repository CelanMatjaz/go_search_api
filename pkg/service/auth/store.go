package auth

import (
	"database/sql"
	"errors"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(connection *db.DbConnection) *Store {
	return &Store{db: connection.DB}
}

func (s *Store) GetUserById(id int) (types.User, error) {
	return s.getUser("WHERE users.id = $1", id)
}

func (s *Store) GetUserByEmail(email string) (types.User, error) {
	return s.getUser("WHERE users.email = $1", email)
}

func (s *Store) CreateUser(user types.User) (types.User, error) {
	transaction, err := s.db.Begin()
	if err != nil {
		return types.User{}, err
	}
	defer transaction.Rollback()

	row := transaction.QueryRow(`
        INSERT INTO users (first_name, last_name, email, password_hash)
        VALUES ($1, $2, $3, $4)
        RETURNING id, first_name, last_name, email, password_hash, created_at, updated_at`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
	)

	newUser, err := scanUserRow(row)
	if err != nil {
		return types.User{}, err
	}

	err = transaction.Commit()
	if err != nil {
		return types.User{}, err
	}

	return newUser, nil
}

func (s *Store) getUser(whereClause string, value any) (types.User, error) {
	row := s.db.QueryRow(`
        SELECT
            id,
            first_name,
            last_name,
            email,
            password_hash,
            created_at,
            updated_at
        FROM users `+whereClause,
		value)

	user, err := scanUserRow(row)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func scanUserRow(row *sql.Row) (types.User, error) {
	var user types.User
	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return types.User{}, types.UserDoesNotExistErr
	}
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}
