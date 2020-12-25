package app

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// https://golang.org/pkg/net/http/#pkg-constants

// The ServerError writes an error message and stack trace to the logger
// and then sends a generic 500 Internal Server Error response to the user.
func (app *App) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Error().Msg(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The ClientError sends a specific status code and corresponding description
// to the user. We'll use this to send responses like 400 "Bad Request" when
// there's a problem with the request that the user sent.
func (app *App) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// This is simply a convenience wrapper around ClientError which sends a
// 404 Not Found response the user.
func (app *App) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
