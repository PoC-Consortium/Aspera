package main

import (
	"flag"
	//	pb "github.com/ac0v/ce-stock/internal/api/protobuf-spec"
	//"fmt"
	"github.com/ac0v/aspera/pkg/p2p"
	"github.com/dgraph-io/badger"
	"github.com/golang/protobuf/proto"
	// "gopkg.in/resty.v1"
	"encoding/binary"
	"log"
	"path/filepath"
	"strconv"
)

func main() {
	path := flag.String("path", "var", "Database Path")
	flag.Parse()

	blockOpts := badger.DefaultOptions
	blockOpts.Dir = filepath.Join(*path, "block")
	blockOpts.ValueDir = blockOpts.Dir
	blockDb, err := badger.Open(blockOpts)
	if err != nil {
		log.Fatal(err)
	}

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	transactionOpts := badger.DefaultOptions
	transactionOpts.Dir = filepath.Join(*path, "transaction")
	transactionOpts.ValueDir = transactionOpts.Dir
	transactionDb, err := badger.Open(transactionOpts)
	if err != nil {
		log.Fatal(err)
	}

	nextBlock := uint64(3444294670862540038)

	for nextBlock != 0 {
		res, _ := p2p.GetNextBlocks(nextBlock)
		for _, block := range res.NextBlocks {
			err = blockDb.Update(func(txn *badger.Txn) error {
				data, _ := proto.Marshal(block)
				key := make([]byte, 8)
				binary.LittleEndian.PutUint64(key, uint64(block.Block))

				err := txn.Set(key, data)
				return err
			})
			for _, transaction := range block.Transactions {
				err = transactionDb.Update(func(txn *badger.Txn) error {
					data, _ := proto.Marshal(transaction)
					key := make([]byte, 8)
					binary.LittleEndian.PutUint64(key, uint64(transaction.Transaction))

					err := txn.Set(key, data)
					return err
				})
			}
			nextBlock, _ = strconv.ParseUint(block.PreviousBlock, 10, 64)
		}
		if len(res.NextBlocks) <= 1 {
			nextBlock = 0
		}
	}
	/*
	   txn := db.NewTransaction(true)
	   for k,v := range updates {
	     if err := txn.Set([]byte(k),[]byte(v)); err == ErrTxnTooBig {
	       _ = txn.Commit()
	       txn = db.NewTransaction(..)
	       _ = txn.Set([]byte(k),[]byte(v))
	     }
	   }
	   _ = txn.Commit()
	*/
	defer transactionDb.Close()
	defer blockDb.Close()
}
