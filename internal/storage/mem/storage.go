package mem

import (
	"sync"

	"github.com/Nexadis/Storage/internal/storage"
)

type MemStorage struct {
	sync.RWMutex
	m map[string]string
}

func (m *MemStorage) Delete(key string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.m, key)
	return nil
}

func (m *MemStorage) Get(key string) (string, error) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.m[key]
	if !ok {
		return "", storage.ErrorNoSuchKey
	}
	return val, nil
}

func (m *MemStorage) Put(key, value string) error {
	m.Lock()
	defer m.Unlock()
	m.m[key] = value

	return nil
}

func New(size int) storage.Storage {
	m := make(map[string]string, size)
	s := &MemStorage{
		m: m,
	}
	return s
}
