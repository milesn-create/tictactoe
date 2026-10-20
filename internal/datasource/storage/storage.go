package storage

import "sync"

type Storage struct {
	Games sync.Map
}

func NewStorage() *Storage {
	return &Storage{}

}
