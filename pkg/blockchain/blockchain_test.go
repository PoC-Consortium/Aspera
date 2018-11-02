package blockchain

import (
	"os"
	"testing"

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
	bc = NewBlockchain().(*blockchain)
	m.Run()
}

func TestNewBlockchain(t *testing.T) {
	assert.NotNil(t, bc.accountDB)
	assert.NotNil(t, bc.blockDB)
}

func TestGetOrCreateAccount(t *testing.T) {
	a, err := bc.GetOrCreateAccount(1337)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(1337), a.Id)
		assert.Equal(t, int64(0), a.Balance)
		assert.Equal(t, nil, a.PublicKey)
	}

	a, err = bc.GetOrCreateAccount(1337)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(1337), a.Id)
		assert.Equal(t, int64(0), a.Balance)
		assert.Equal(t, nil, a.PublicKey)
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	type test struct {
		id     uint64
		amount int64
		err    error
	}

	tests := []test{
		test{id: 1, amount: 1, err: nil},
		test{id: 1, amount: -1, err: nil},
		test{id: 1, amount: -1, err: ErrBalanceTooLow},
		test{id: 2, amount: -1, err: ErrBalanceTooLow},
	}

	for _, test := range tests {
		err := bc.UpdateAccountBalance(test.id, test.amount)
		assert.Equal(t, test.err, err)
	}
}
