package store

import (
	"github.com/PoC-Consortium/Aspera/pkg/config"
)

type Store struct {
	RawStore *RawStore
}

func Init(path string, genesisMilestone config.Milestone) *Store {
	var store Store
	store.RawStore = NewRawStore(path, genesisMilestone)
	return &store
}

func (store *Store) Close() {
}
