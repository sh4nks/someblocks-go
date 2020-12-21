package app

import (
	"context"

	"github.com/alexedwards/scs/v2"
)

type Flash struct {
	Level   string
	Message string
}

type SessionManager struct {
	scs.SessionManager
}

func (s *SessionManager) SetFlash(ctx context.Context, level string, message string) {
	flash := Flash{
		Level:   level,
		Message: message,
	}
	s.Put(ctx, "flashed_messages", flash)
}

func (s *SessionManager) GetFlash(ctx context.Context) *Flash {
	x := s.Pop(ctx, "flashed_messages")
	flash, ok := x.(Flash)
	if !ok {
		return nil
	}
	return &flash
}
