package forms

import (
	"net/url"
	"someblocks/models"
)

type ProfileForm struct {
	Form
}

func NewProfileForm(data url.Values) *RegisterForm {
	return &RegisterForm{
		Form{
			data,
			errors(map[string][]string{}),
		},
	}
}

func (f *ProfileForm) Valid(us *models.UserService) bool {
	// all info is optional at the moment
	return len(f.Errors) == 0
}
