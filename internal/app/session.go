package app

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type Flash struct {
	Level   string
	Message string
}

type SessionManager struct {
	scs.SessionManager
}

func (s *SessionManager) setFlash(r *http.Request, message string, level string) {
	flash := Flash{
		Level:   level,
		Message: message,
	}
	s.Put(r.Context(), "flashed_messages", flash)
}

func (s *SessionManager) getFlash(r *http.Request) *Flash {
	x := s.Pop(r.Context(), "flashed_messages")
	flash, ok := x.(Flash)
	if !ok {
		return nil
	}
	return &flash
}

func (app *App) Flash(r *http.Request, message string, level string) {
	app.Session.setFlash(r, message, level)
}
