package models

import (
	"strings"
)

const (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound modelError = "models: Resource not found"
	// ErrInvalidCredentials is returned if a user tries to login with an
	// incorrect email address or password.
	ErrInvalidLoginCredentials modelError = "models: Invalid login credentials."
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
