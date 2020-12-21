package app

import (
	"context"
	"fmt"

	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog/log"
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
	log.Debug().Msgf("SET: %v", flash)
	s.Put(ctx, "flashed_messages", flash)
}

func (s *SessionManager) GetFlash(ctx context.Context) *Flash {
	x := s.Pop(ctx, "flashed_messages")
	log.Debug().Msgf("GET: %v", x)
	flash, ok := x.(Flash)
	if !ok {
		fmt.Printf("Couldn't convert x %v to %v", x, flash)
		return nil
	}
	return &flash
}
