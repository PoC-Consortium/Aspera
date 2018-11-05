package blockchain

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/ac0v/aspera/pkg/account"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
)

var bc *blockchain

func TestMain(m *testing.M) {
	if err := os.RemoveAll("var"); err != nil {
		panic(err)
	}
	if err := os.Mkdir("var", 0777); err != nil {
		panic(err)
	}
	Init()
	bc = BC.(*blockchain)
	m.Run()
}

func createDummyAccount(a *account.Account) {
	err := bc.accountDB.Update(func(txn *badger.Txn) error {
		if err := bc.accountIndex.Index(a); err != nil {
			return err
		}
		return bc.updateAccount(txn, a)
	})
	if err != nil {
		panic(err)
	}
}

func TestNewBlockchain(t *testing.T) {
	assert.NotNil(t, bc.accountDB)
	assert.NotNil(t, bc.accountIndex)
	assert.NotNil(t, bc.blockDB)
}

func TestSendBurst(t *testing.T) {
	senderPublicKey, _ := hex.DecodeString("d37ebb299c54fa4603f1b656f6bcf70810a9cd1d56e2ea979c7933d94ae3602a")
	var receiverId uint64 = 8964081770308214857

	err := bc.SendBurst(senderPublicKey, receiverId, 10, 1)
	assert.Equal(t, ErrUnknownAccount, err)

	sender := account.NewAccount(7900104405094198526)
	sender.Balance = 10
	sender.PublicKey = senderPublicKey
	createDummyAccount(sender)

	err = bc.SendBurst(senderPublicKey, receiverId, 10, 1)
	assert.Equal(t, ErrBalanceTooLow, err)

	sender.Balance = 11
	createDummyAccount(sender)

	err = bc.SendBurst(senderPublicKey, receiverId, 10, 1)
	assert.Nil(t, err)

	receiver, err := bc.GetAccountById(receiverId)
	if assert.Nil(t, err) && assert.NotNil(t, receiver) {
		assert.Equal(t, receiverId, receiver.Id)
		assert.Nil(t, receiver.PublicKey)
		assert.Equal(t, int64(10), receiver.Balance)
		// TODO: when to set address and test!
	}

	sender, err = bc.GetAccountByAddress("8V9Y-58B4-RVWP-8HQAV")
	// sender, err = bc.GetAccountByPublicKey(senderPublicKey)
	if assert.Nil(t, err) && assert.NotNil(t, sender) {
		assert.Equal(t, senderPublicKey, sender.PublicKey)
		assert.Equal(t, int64(0), sender.Balance)
	}
}
