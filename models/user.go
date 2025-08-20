package models

import (
	"database/sql"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	ID        int    `gorm:"primarykey"`
	Email     string `gorm:"unique;size:254"` // emails should not exceed 254 characters
	Username  string `gorm:"unique;size:128"`
	Password  []byte `gorm:"size:256"`
	Role      string `gorm:"size:128"`
	LastLogin sql.NullTime
}

type UserService struct {
	DB *gorm.DB
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

func (u *UserService) Insert(username, email, password string) (*User, error) {
	// 2^12 = 4096
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     "user",
	}

	result := u.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *UserService) Update(user *User) error {
	result := u.DB.Model(&user).Select("*").Omit("password", "role").Updates(User{Username: user.Username, Email: user.Email})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserService) UpdatePassword(user *User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	// result := u.DB.Select("password").Updates(user)
	result := u.DB.Model(&user).Select("password").Where("username = ?", user.Username).Updates(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserService) UpdateRole(user *User, role string) error {
	user.Role = role
	result := u.DB.Model(&user).Select("role").Where("username = ?", user.Username).Updates(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserService) Search(username string, resultSize int) []User {
	var result *gorm.DB
	var users []User
	if len(username) > 0 {
		log.Debug().Msgf("searching for 'username' LIKE ", username)
		result = u.DB.Where("username LIKE ?", username).Find(&users)
	} else {
		log.Debug().Msgf("searching for all users")
		result = u.DB.Find(&users)
	}

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("UserService.Search:")
	}
	return users
}

// GetById returns the User struct or nil if no user matching the id has not
// been found.
func (u *UserService) GetById(id int) *User {
	user := User{}
	result := u.DB.First(&user, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("UserService.GetById:")
		return nil
	}
	return &user
}

func (u *UserService) GetByUsername(username string) *User {
	user := User{}
	result := u.DB.Where("username = ?", username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("UserService.GetByUsername:")
		return nil
	}
	return &user
}

func (u *UserService) GetByEmail(email string) *User {
	user := User{}
	result := u.DB.Where("email = ?", email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("UserService.GetByEmail:")
		return nil
	}
	return &user
}
