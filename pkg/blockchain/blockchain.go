package blockchain

import (
	"errors"

	"github.com/ac0v/aspera/pkg/account"
	"github.com/ac0v/aspera/pkg/blockchain/index"
	"github.com/ac0v/aspera/pkg/crypto"

	"github.com/dgraph-io/badger"
)

var BC Blockchain

var (
	ErrBalanceTooLow  = errors.New("account's balance too low")
	ErrUnknownAccount = errors.New("account unknown")
)

type Blockchain interface {
	GetAccountById(id uint64) (*account.Account, error)
	GetAccountByPublicKey(publicKey []byte) (*account.Account, error)
	GetAccountByAddress(address string) (*account.Account, error)
	SendBurst(senderPublicKey []byte, receiverId uint64, amount, fee int64) error
}

type blockchain struct {
	accountDB    *badger.DB
	accountIndex index.AccountIndex

	blockDB *badger.DB
}

type accountIndexData struct {
	PublicKey string
	Address   string
}

func Init() {
	BC = newBlockchain()
}

func openDB(name string) *badger.DB {
	opts := badger.DefaultOptions
	opts.Dir = "var/" + name
	opts.ValueDir = "var/" + name
	if db, err := badger.Open(opts); err == nil {
		return db
	} else {
		panic(err)
	}
}

func newBlockchain() Blockchain {
	return &blockchain{
		accountDB:    openDB("account"),
		accountIndex: index.NewAccountIndex(),

		blockDB: openDB("block"),
	}
}

func (bc *blockchain) updateAccount(txn *badger.Txn, a *account.Account) error {
	return txn.Set(bc.accountIndex.ById(a.Id), a.ToBytes())
}

func (bc *blockchain) getAccount(txn *badger.Txn, key []byte) (*account.Account, error) {
	if v, err := txn.Get(key); err == nil {
		if bs, err := v.Value(); err == nil {
			return account.FromBytes(bs), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (bc *blockchain) getAccountById(txn *badger.Txn, id uint64) (*account.Account, error) {
	key := bc.accountIndex.ById(id)
	return bc.getAccount(txn, key)
}

func (bc *blockchain) getAccountByPublicKey(txn *badger.Txn, pubKey []byte) (*account.Account, error) {
	if key, err := bc.accountIndex.ByPublicKey(pubKey); err == nil {
		return bc.getAccount(txn, key)
	} else {
		return nil, err
	}
}

func (bc *blockchain) getAccountByAddress(txn *badger.Txn, address string) (*account.Account, error) {
	if key, err := bc.accountIndex.ByAddress(address); err == nil {
		return bc.getAccount(txn, key)
	} else {
		return nil, err
	}
}

func (bc *blockchain) GetAccountById(id uint64) (*account.Account, error) {
	var a *account.Account
	err := bc.accountDB.View(func(txn *badger.Txn) error {
		var err error
		a, err = bc.getAccountById(txn, id)
		return err
	})
	return a, err
}

func (bc *blockchain) GetAccountByPublicKey(publicKey []byte) (*account.Account, error) {
	var a *account.Account
	err := bc.accountDB.View(func(txn *badger.Txn) error {
		var err error
		a, err = bc.getAccountByPublicKey(txn, publicKey)
		return err
	})
	return a, err
}

func (bc *blockchain) GetAccountByAddress(address string) (*account.Account, error) {
	var a *account.Account
	err := bc.accountDB.View(func(txn *badger.Txn) error {
		var err error
		a, err = bc.getAccountByAddress(txn, address)
		return err
	})
	return a, err
}

func (bc *blockchain) SendBurst(senderPublicKey []byte, receiverId uint64, amount, fee int64) error {
	return bc.accountDB.Update(func(txn *badger.Txn) error {
		sender, err := bc.getAccountByPublicKey(txn, senderPublicKey)
		switch err {
		case nil:
		case badger.ErrKeyNotFound:
			// an account only gets a public key on its first outgoing transaction
			// so we need to check if we can find it by numeric id
			_, id := crypto.BytesToHashAndID(senderPublicKey)
			switch sender, err = bc.getAccountById(txn, id); err {
			case nil:
				// we found the account by numeric id, but not by public key
				// this is the account's first outgoing transaction
				// we can no activate it by setting its public key
				sender.PublicKey = senderPublicKey

				// to be retrievable by public key we need to add it to
				// the account index
				if err := bc.accountIndex.Index(sender); err != nil {
					return err
				}
			case badger.ErrKeyNotFound:
				return ErrUnknownAccount
			default:
				return err
			}
		default:
			return err
		}

		totalAmount := amount + fee
		if sender.Balance < totalAmount {
			return ErrBalanceTooLow
		}
		sender.Balance -= totalAmount

		receiver, err := bc.getAccountById(txn, receiverId)
		switch err {
		case nil:
		case badger.ErrKeyNotFound:
			// receiver did not exist yet, so its a new account that still
			// needs to get activated with an outgoing transactio
			receiver = account.NewAccount(receiverId)
		default:
			return err
		}
		receiver.Balance += amount

		if err := bc.updateAccount(txn, sender); err != nil {
			return err
		}
		return bc.updateAccount(txn, receiver)
	})
}
