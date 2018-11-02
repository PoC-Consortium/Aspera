package blockchain

import (
	"encoding/hex"
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

func TestUpdateAccountBalance(t *testing.T) {
	type test struct {
		id        uint64
		publicKey string
		amount    int64
		err       error
	}

	tests := []test{
		test{id: 7900104405094198526, amount: 1, publicKey: "d37ebb299c54fa4603f1b656f6bcf70810a9cd1d56e2ea979c7933d94ae3602a", err: nil},
		test{id: 7900104405094198526, amount: -1, publicKey: "d37ebb299c54fa4603f1b656f6bcf70810a9cd1d56e2ea979c7933d94ae3602a", err: nil},
		test{id: 7900104405094198526, amount: -1, publicKey: "d37ebb299c54fa4603f1b656f6bcf70810a9cd1d56e2ea979c7933d94ae3602a", err: ErrBalanceTooLow},
		test{id: 8964081770308214857, amount: -1, publicKey: "22c1175730c735c7b768f32b25e426454658ec2cb58f49723412ed49e9cf5a44", err: ErrBalanceTooLow},
	}

	for _, test := range tests {
		publiyKey, _ := hex.DecodeString(test.publicKey)
		err := bc.UpdateAccountBalance(test.id, test.amount, publiyKey)
		assert.Equal(t, test.err, err)
	}
}
