package storage

import "errors"

type Storage interface {
	Save(p *Page) error
	PickRandom(username string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	ID       int64
	URL      string
	UserName string
}

var (
	ErrNotFound = errors.New("pages not found")
)
