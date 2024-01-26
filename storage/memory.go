package storage

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
)

type memoryStorage struct {
	storage map[uuid.UUID]File
	mutex   sync.Mutex
}

type memoryFile struct {
	content   []byte
	mime      string
	expiredAt *time.Time
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		storage: make(map[uuid.UUID]File),
	}
}

func (f *memoryFile) GetContent() []byte {
	return f.content
}

func (f *memoryFile) GetMime() string {
	return f.mime
}

func (f *memoryFile) GetExpiredAt() *time.Time {
	return f.expiredAt
}

func (s *memoryStorage) Save(file File) uuid.UUID {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := uuid.New()
	s.storage[id] = file
	return id
}

func (s *memoryStorage) Get(id uuid.UUID) (File, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, ok := s.storage[id]
	return file, ok
}

func (s *memoryStorage) GC() {
	// fixme: stop-the-world
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, file := range s.storage {
		if file.GetExpiredAt() == nil {
			continue
		}

		if IsExpired(file) {
			delete(s.storage, id)
		}
	}
}

func (f *memoryFile) GetContentReader() io.Reader {
	return bytes.NewReader(f.content)
}

var _ Storage = (*memoryStorage)(nil)
var _ File = (*memoryFile)(nil)
