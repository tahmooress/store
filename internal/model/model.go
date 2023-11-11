package model

import (
	"io"
	"time"
)

type FileObject struct {
	ID              int64
	OwnerID         string
	Name            string
	Type            string
	StorageLocation string
	Content         io.Reader
	Size            int64
	Tags            []string
	CreatedAt       time.Time
}

type Filter struct {
	ownerID string
	name    string
	tags    []string
}

func NewFilter(ownerID string) *Filter {
	return &Filter{
		ownerID: ownerID,
	}
}

func (f *Filter) SetName(name string) *Filter {
	f.name = name
	return f
}

func (f *Filter) SetTags(tags []string) *Filter {
	f.tags = make([]string, len(tags))
	copy(f.tags, tags)
	return f
}

type FilterQueryBuilder func(owner_id, name string, tags []string) (
	query string, args []interface{}, err error)

func (f *Filter) Query(fn FilterQueryBuilder) (string, []interface{}, error) {
	return fn(f.ownerID, f.name, f.tags)
}
