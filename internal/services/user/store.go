package user

import (
	"errors"
	"finance-crud-app/internal/types"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

var (
	CreateUserError   = errors.New("cannot create user")
	RetrieveUserError = errors.New("cannot retrieve user")
	DeleteUserError   = errors.New("cannot delete user")
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) (user_id int, err error) {
	query := `
	INSERT INTO users
	(firstName, lastName, email, password)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	var userId int
	err = s.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&userId)
	if err != nil {
		return -1, CreateUserError
	}

	return userId, nil
}

func (s *Store) GetUserByEmail(email string) (types.User, error) {
	var user types.User

	err := s.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return types.User{}, RetrieveUserError
	}

	if user.ID == 0 {
		log.Fatalf("user not found")
		return types.User{}, RetrieveUserError
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	var user types.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, RetrieveUserError
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (s *Store) DeleteUser(email string) error {

	user, err := s.GetUserByEmail(email)
	if err != nil {
		return DeleteUserError
	}
	// delete user records first
	_, err = s.db.Exec("DELETE FROM records WHERE userid = $1", user.ID)
	if err != nil {
		return DeleteUserError
	}

	_, err = s.db.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		return DeleteUserError
	}
	return nil
}
