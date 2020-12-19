package forms

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type NewSnippet struct {
	Title    string
	Content  string
	Expires  string
	Failures map[string]string
}

type SignupUser struct {
	Name     string
	Email    string
	Password string
	Failures map[string]string
}

type LoginUser struct {
	Email    string
	Password string
	Failures map[string]string
}

func (ns *NewSnippet) Valid() bool {
	ns.Failures = make(map[string]string)

	if len(strings.TrimSpace(ns.Title)) == 0 {
		ns.Failures["Title"] = "Title is required"
	} else if len(strings.TrimSpace(ns.Title)) > 50 {
		ns.Failures["Title"] = "Title cannot be longer than 50 characters"
	}

	if len(strings.TrimSpace(ns.Content)) == 0 {
		ns.Failures["Content"] = "Content is required"
	}

	permitted := map[string]bool{"1209600": true, "604800": true, "86400": true}
	if len(strings.TrimSpace(ns.Expires)) == 0 {
		ns.Failures["Expires"] = "Expiry time is required"
	} else if !permitted[ns.Expires] {
		ns.Failures["Expires"] = "Expiry time must be 1209600, 604800 or 86400 seconds"
	}

	return len(ns.Failures) == 0
}

func (su *SignupUser) Valid() bool {
	su.Failures = make(map[string]string)

	if len(strings.TrimSpace(su.Name)) == 0 {
		su.Failures["Name"] = "Name is required"
	} else if len(strings.TrimSpace(su.Name)) > 50 || len(strings.TrimSpace(su.Name)) < 3 {
		su.Failures["Name"] = "Name cannot be longer than 50 and shorter than 3 characters"
	}

	if len(strings.TrimSpace(su.Email)) == 0 {
		su.Failures["Email"] = "Email is required"
	} else if len(strings.TrimSpace(su.Email)) > 254 || !rxEmail.MatchString(su.Email) {
		su.Failures["Email"] = "Email is not valid"
	}

	if len(strings.TrimSpace(su.Password)) == 0 {
		su.Failures["Password"] = "Password is required"
	} else if utf8.RuneCountInString(su.Password) < 8 {
		su.Failures["Password"] = "Password is too short"
	}

	return len(su.Failures) == 0
}

func (lu *LoginUser) Valid() bool {
	lu.Failures = make(map[string]string)

	if len(strings.TrimSpace(lu.Email)) == 0 {
		lu.Failures["Email"] = "Email is required"
	} else if len(strings.TrimSpace(lu.Email)) > 254 || !rxEmail.MatchString(lu.Email) {
		lu.Failures["Email"] = "Email is not valid"
	}

	if len(strings.TrimSpace(lu.Password)) == 0 {
		lu.Failures["Password"] = "Password is required"
	} else if utf8.RuneCountInString(lu.Password) < 8 {
		lu.Failures["Password"] = "Password is too short"
	}

	return len(lu.Failures) == 0
}
