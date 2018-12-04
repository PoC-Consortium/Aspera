package blockchain

import (
	"errors"

	"github.com/PoC-Consortium/aspera/pkg/account"
	"github.com/PoC-Consortium/aspera/pkg/crypto"

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
	TransferBurst(senderPublicKey []byte, receiverId uint64, amount, fee int64) error
}

type blockchain struct {
	// one or multiple badger dbs?
	// + multiple: more key flexibility
	// - multiple: no transactions between databases
	db *badger.DB
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
		db: openDB("blockchain"),
	}
}

func (bc *blockchain) indexAccountsPublicKey(txn *badger.Txn, a *account.Account) error {
	return txn.Set(a.PublicKey, uint64ToBs(a.Id))
}

func (bc *blockchain) updateAccount(txn *badger.Txn, a *account.Account) error {
	return txn.Set(uint64ToBs(a.Id), a.ToBytes())
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
	key := uint64ToBs(id)
	return bc.getAccount(txn, key)
}

func (bc *blockchain) getAccountByPublicKey(txn *badger.Txn, pubKey []byte) (*account.Account, error) {
	if v, err := txn.Get(pubKey); err == nil {
		if bs, err := v.Value(); err == nil {
			return bc.getAccount(txn, bs)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (bc *blockchain) GetAccountById(id uint64) (*account.Account, error) {
	var a *account.Account
	err := bc.db.View(func(txn *badger.Txn) error {
		var err error
		a, err = bc.getAccountById(txn, id)
		return err
	})
	return a, err
}

func (bc *blockchain) GetAccountByPublicKey(publicKey []byte) (*account.Account, error) {
	var a *account.Account
	err := bc.db.View(func(txn *badger.Txn) error {
		var err error
		a, err = bc.getAccountByPublicKey(txn, publicKey)
		return err
	})
	return a, err
}

// TransferBurst transfers burst from one account to another.
// If sender's public key is nil only money on the reciver sice will be added/created (block forging).
func (bc *blockchain) TransferBurst(senderPublicKey []byte, receiverId uint64, amount, fee int64) error {
	totalAmount := amount + fee
	return bc.db.Update(func(txn *badger.Txn) error {
		// senderPublicKey nil -> block forge
		if senderPublicKey != nil {
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
					if err := bc.indexAccountsPublicKey(txn, sender); err != nil {
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

			if sender.Balance < totalAmount {
				return ErrBalanceTooLow
			}
			sender.Balance -= totalAmount
			if err := bc.updateAccount(txn, sender); err != nil {
				return err
			}
		}

		receiver, err := bc.getAccountById(txn, receiverId)
		switch err {
		case nil:
			// if this is a block forge transfer -> reward recipient handling
			if senderPublicKey == nil && receiver.RewardRecipient != receiver.Id {
				receiver, err = bc.getAccountById(txn, receiverId)
				switch err {
				case nil:
				case badger.ErrKeyNotFound:
					receiver = account.NewAccount(receiverId)
				default:
					return err
				}
			}
		case badger.ErrKeyNotFound:
			// receiver did not exist yet, so its a new account that still
			// needs to get activated with an outgoing transaction
			receiver = account.NewAccount(receiverId)
		default:
			return err
		}
		receiver.Balance += amount

		return bc.updateAccount(txn, receiver)
	})
}
