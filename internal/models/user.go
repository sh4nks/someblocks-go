package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int          `db:"id"`
	Email      string       `db:"email"`
	Username   string       `db:"username"`
	Password   []byte       `db:"password"`
	LastLogin  sql.NullTime `db:"last_login"`
	CreatedAt  time.Time    `db:"created_at"`
	ModifiedAt time.Time    `db:"modified_at"`
}

type UserService struct {
	DB *sqlx.DB
}

// Authenticate is used to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (u *UserService) Authenticate(email, password string) (int, error) {
	user := u.GetByEmail(email)
	if user == nil {
		return 0, ErrInvalidLoginCredentials
	}

	// Check whether the hashed password and plain-text password provided match
	// If they don't, we return the ErrInvalidLoginCredentials error.
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidLoginCredentials
	} else if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (u *UserService) Insert(username, email, password string) error {
	// 2^12 = 4096
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = u.DB.NamedExec(`
		INSERT INTO users (
			username, email, password
		) VALUES (
			:username, :email, :password
		)`,
		User{
			Username: username,
			Email:    email,
			Password: hashedPassword,
		})
	return err
}

// GetById returns the User struct or nil if no user matching the id has not
// been found.
func (u *UserService) GetById(id int) *User {
	user := User{}
	err := u.DB.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Error().Err(err).Msg("UserService.GetById:")
		return nil
	}
	return &user
}

func (u *UserService) GetByUsername(username string) *User {
	user := User{}
	err := u.DB.Get(&user, "SELECT * FROM users WHERE username = ?", username)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Error().Err(err).Msg("UserService.GetByUsername:")
		return nil
	}
	return &user
}

func (u *UserService) GetByEmail(email string) *User {
	user := User{}
	err := u.DB.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Error().Err(err).Msg("UserService.GetByEmail:")
		return nil
	}
	return &user
}
