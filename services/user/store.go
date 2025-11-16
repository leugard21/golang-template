package user

import (
	"database/sql"
	"template/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) (int, error) {
	var userID int

	err := s.db.QueryRow(
		"INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Username,
		user.Email,
		user.Password,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, email, password, role, created_at
         FROM users
         WHERE email = $1
         LIMIT 1`,
		email,
	)

	var u types.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Store) GetUserByUsername(username string) (*types.User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, email, password, role, created_at
         FROM users
         WHERE username = $1
         LIMIT 1`,
		username,
	)

	var u types.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, email, password, role, created_at
         FROM users
         WHERE id = $1
         LIMIT 1`,
		id,
	)

	var u types.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
