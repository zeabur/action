// Package storage provides a memory-based storage for files.
// Note that it is for a short time.
package storage

import (
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	Save(file File) uuid.UUID
	Get(id uuid.UUID) (File, bool)
	GC()
}

type File interface {
	GetContent() []byte
	GetMime() string
	GetExpiredAt() *time.Time
}

func IsExpired(f File) bool {
	if f.GetExpiredAt() == nil {
		return false
	}

	return f.GetExpiredAt().Before(time.Now())
}
