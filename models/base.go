package models

import (
	"time"
)

type BaseModel struct {
	ModelTimestamps
}

type ModelTimestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
