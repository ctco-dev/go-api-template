package beer

import (
	"context"
)

// ID is string alias type for User id
type ID = string

// Beer represents a beer
type Beer struct {
	ID   ID
	Name string
}

// Reader reads beer by id
type Reader interface {
	Read(context.Context, ID) (*Beer, error)
}

// Writer writes a new beer and returns its id
type Writer interface {
	Write(context.Context, Beer) (ID, error)
}

// Remover removes a beer by id
type Remover interface {
	Remove(context.Context, ID) error
}

// Repo is an interface for a beer repository
type Repository interface {
	Reader
	Writer
	Remover
}
