package forms

import (
	"net/url"
	"someblocks/internal/models"
)

type LoginForm struct {
	Form
}

func NewLoginForm(data url.Values) *LoginForm {
	return &LoginForm{
		Form{
			data,
			errors(map[string][]string{}),
		},
	}
}

func (f *LoginForm) Valid() bool {
	f.Required("email", "password")
	return len(f.Errors) == 0
}

type RegisterForm struct {
	Form
}

func NewRegisterForm(data url.Values) *RegisterForm {
	return &RegisterForm{
		Form{
			data,
			errors(map[string][]string{}),
		},
	}
}

func (f *RegisterForm) Valid(us *models.UserService) bool {
	f.Required("username")
	f.MinLength("username", 3)

	f.Required("email")
	f.ValidEmail("email")

	f.Required("password")
	f.MinLength("password", 8)

	f.Required("confirm_password")
	f.MinLength("confirm_password", 8)

	f.ConfirmPassword("password", "confirm_password")

	f.Required("tos")

	value := f.Get("username")
	if value != "" {
		if user := us.GetByUsername(value); user != nil {
			f.Errors.Add("username", "Username already in use.")
		}
	}

	value = f.Get("email")
	if value != "" {
		if user := us.GetByEmail(value); user != nil {
			f.Errors.Add("email", "Email address already in use.")
		}
	}

	return len(f.Errors) == 0
}
