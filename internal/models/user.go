package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Username  string    `db:"username"`
	Password  []byte    `db:"password"`
	LastLogin time.Time `db:"last_login"`
	CreatedAt time.Time `db:"created_at"`
}

type UserService struct {
	DB *sqlx.DB
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

func (u *UserService) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (u *UserService) GetById(id int) (*User, error) {
	return nil, nil
}
