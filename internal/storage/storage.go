package storage

import (
	"errors"
)

const DefaultUser = "default"

var ErrorNoSuchKey = errors.New("no such key")

type Storage interface {
	Put(user, key, value string) error
	Get(user, key string) (string, error)
	Delete(user, key string) error
}
