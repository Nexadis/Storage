package mem

import (
	"sync"

	"github.com/Nexadis/Storage/internal/storage"
)

type MemStorage struct {
	sync.RWMutex
	m map[string]*userStorage
}

type userStorage struct {
	sync.RWMutex
	m map[string]string
}

func newUserStorage() *userStorage {
	return &userStorage{
		m: make(map[string]string),
	}
}

func (m *MemStorage) Delete(user, key string) error {
	m.Lock()
	kv, ok := m.m[user]
	m.Unlock()
	if !ok {
		return nil
	}
	delete(kv.m, key)
	return nil
}

func (m *MemStorage) Get(user, key string) (string, error) {
	m.RLock()
	kv, ok := m.m[user]
	m.RUnlock()
	if !ok {
		return "", storage.ErrorNoSuchKey
	}
	kv.RLock()
	val, ok := kv.m[key]
	kv.RUnlock()
	if !ok {
		return "", storage.ErrorNoSuchKey
	}
	return val, nil
}

func (m *MemStorage) Put(user, key, value string) error {
	m.RLock()
	kv, ok := m.m[user]
	m.RUnlock()
	if !ok {
		kv = newUserStorage()
		m.Lock()
		m.m[user] = kv
		m.Unlock()
	}
	kv.Lock()
	kv.m[key] = value
	kv.Unlock()

	return nil
}

func New(size int) storage.Storage {
	m := make(map[string]*userStorage, size)
	s := &MemStorage{
		m: m,
	}
	return s
}
