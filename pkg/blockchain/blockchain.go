package blockchain

import (
	"encoding/binary"
	"errors"

	"github.com/ac0v/aspera/pkg/account"

	"github.com/dgraph-io/badger"
)

var (
	ErrBalanceTooLow = errors.New("account's balance too low")
)

type Blockchain interface {
	GetAccount(id uint64) (*account.Account, error)
	UpdateAccountBalance(id uint64, amount int64, publicKey []byte) error
}

type blockchain struct {
	accountDB *badger.DB
	blockDB   *badger.DB
}

func db(name string) *badger.DB {
	opts := badger.DefaultOptions
	opts.Dir = "var/" + name
	opts.ValueDir = "var/" + name
	if db, err := badger.Open(opts); err == nil {
		return db
	} else {
		panic(err)
	}
}

func NewBlockchain() Blockchain {
	return &blockchain{
		accountDB: db("account"),
		blockDB:   db("block"),
	}
}

func accountIDToBytes(id uint64) []byte {
	var idBs [8]byte
	binary.LittleEndian.PutUint64(idBs[:], id)
	return idBs[:]
}

func (bc *blockchain) updateAccount(txn *badger.Txn, a *account.Account) error {
	return txn.Set(accountIDToBytes(a.Id), a.ToBytes())
}

func (bc *blockchain) getAccount(txn *badger.Txn, id uint64) (*account.Account, error) {
	if v, err := txn.Get(accountIDToBytes(id)); err == nil {
		if bs, err := v.Value(); err == nil {
			return account.FromBytes(bs), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (bc *blockchain) GetAccount(id uint64) (*account.Account, error) {
	var a *account.Account
	err := bc.accountDB.View(func(txn *badger.Txn) error {
		var err error
		a, err = bc.getAccount(txn, id)
		return err
	})
	return a, err
}

func (bc *blockchain) UpdateAccountBalance(id uint64, amount int64, publicKey []byte) error {
	err := bc.accountDB.Update(func(txn *badger.Txn) error {
		a, err := bc.getAccount(txn, id)

		switch err {
		case nil:
			a.Balance += amount
			if a.Balance < 0 {
				return ErrBalanceTooLow
			}
			return bc.updateAccount(txn, a)
		case badger.ErrKeyNotFound:
			if amount >= 0 {
				return bc.updateAccount(txn, account.NewAccount(publicKey, amount))
			} else {
				return ErrBalanceTooLow
			}
		default:
			return err
		}
	})
	return err
}
