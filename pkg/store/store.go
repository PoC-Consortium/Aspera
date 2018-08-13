package store

import (
	r "github.com/ac0v/aspera/pkg/registry"
)

type Store struct {
	RawStore *RawStore
}

func Init(registry *r.Registry) *Store {
	var store Store
	store.RawStore = NewRawStore(registry)
	return &store
}

func (store *Store) Close() {
}
