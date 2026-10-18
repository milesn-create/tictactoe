package datasource

import "sync"

type Storage struct {
	games sync.Map
}

func NewStorage() *Storage {
	return &Storage{}

}
