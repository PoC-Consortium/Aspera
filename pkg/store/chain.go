package store

import (
	"os"
	"path/filepath"

	"github.com/PoC-Consortium/Aspera/pkg/block"
	. "github.com/PoC-Consortium/Aspera/pkg/log"

	"github.com/dgraph-io/badger"
	"go.uber.org/zap"
)

type ChainStore struct {
	db *badger.DB
}

func NewChainStore(path string) *ChainStore {
	var chainStore ChainStore

	basePath := filepath.Join(path, "chain")
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, os.ModePerm)
	}

	opts := badger.DefaultOptions
	opts.Dir = basePath
	opts.ValueDir = basePath
	db, err := badger.Open(opts)
	if err != nil {
		zap.Error(err)
	}
	chainStore.db = db

	return &chainStore
}

func (chainStore *ChainStore) Store(block *block.Block) {
	Log.Info("Store Block to ChainStore", zap.Int32("height", block.Height), zap.Uint64("id", block.Id))
	txn := chainStore.db.NewTransaction(true)
	if blockDataBs, err := block.Freeze(); err == nil {
		for _, keyValue := range blockDataBs {
			switch err := txn.Set(keyValue[0], keyValue[1]); err {
			case nil:
			case badger.ErrTxnTooBig:
				if err = txn.Commit(nil); err != nil {
					panic(err)
				}
				txn = chainStore.db.NewTransaction(true)
				if err := txn.Set(keyValue[0], keyValue[1]); err != nil {
					panic(err)
				}
			default:
				panic(err)
			}
		}
	} else {
		panic(err)
	}
	if err := txn.Commit(nil); err != nil {
		panic(err)
	}
}

func (chainStore *ChainStore) FindBlockBy(key string) (*block.Block, error) {
	var b *block.Block
	err := chainStore.db.View(func(txn *badger.Txn) error {
		var err error
		b, err = block.Thaw(txn, key)
		return err
	})
	return b, err
}

func (chainStore *ChainStore) FindBlocksAfter(key string, limit int) ([]*block.Block, error) {
	var blocks []*block.Block
	err := chainStore.db.View(func(txn *badger.Txn) error {
		var err error
		blocks, err = block.BulkThaw(txn, key, limit)
		return err
	})
	return blocks, err
}

func (chainStore *ChainStore) Rebuild() {
	for rit := s.RawStore.Iterator(); rit.Next(); {
		chainStore.Store(rit.Current)
	}
}
