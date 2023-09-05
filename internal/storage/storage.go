package storage

import "errors"

var ErrorNoSuchKey = errors.New("no such key")

type Storage interface {
	Put(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}
