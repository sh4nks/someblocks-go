package forms

import (
	"net/url"
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
	f.Required("Email", "Password")
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

func (f *RegisterForm) Valid() bool {
	f.Required("Email")
	f.ValidEmail("Email")

	f.Required("Password")
	f.MinLength("Password", 8)
	return len(f.Errors) == 0
}
