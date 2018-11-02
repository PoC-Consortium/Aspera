package account

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	type test struct {
		publicKey string
		balance   int64
		id        uint64
		address   string
	}

	tests := []test{
		test{
			publicKey: "d37ebb299c54fa4603f1b656f6bcf70810a9cd1d56e2ea979c7933d94ae3602a",
			balance:   2,
			id:        7900104405094198526,
			address:   "8V9Y-58B4-RVWP-8HQAV",
		},
		test{
			publicKey: "22c1175730c735c7b768f32b25e426454658ec2cb58f49723412ed49e9cf5a44",
			balance:   3,
			id:        8964081770308214857,
			address:   "FA4B-W3EE-UY38-9JQTS",
		},
	}

	for _, test := range tests {
		publicKey, _ := hex.DecodeString(test.publicKey)

		a := NewAccount(publicKey, test.balance)

		assert.Equal(t, publicKey, a.PublicKey)
		assert.Equal(t, test.balance, a.Balance)
		assert.Equal(t, test.id, a.Id)
		assert.Equal(t, test.id, a.RewardRecipient)
		assert.Equal(t, test.address, a.Address)
	}
}
