package store

import (
	"os"
	"path/filepath"

	"github.com/PoC-Consortium/Aspera/pkg/block"

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
	txn := chainStore.db.NewTransaction(true)
	if block, err := block.Thaw(txn, key); err == nil {
		if err := txn.Commit(nil); err != nil {
			panic(err)
		}
		return block, err
	} else {
		return block, err
	}

}
