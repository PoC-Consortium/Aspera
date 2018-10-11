package block

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ac0v/aspera/pkg/transaction"
)

func TestGetID(t *testing.T) {
	for _, blockTest := range BlockTests {
		b := blockTest.Block

		// dummy set, because we need to hash over number of transactions
		b.Transactions = make([]*transaction.Transaction, blockTest.TXLen)

		if blockTest.BlockATs != "" {
			b.BlockATs = &blockTest.BlockATs
		}

		id, err := b.CalculateID()
		if assert.Nil(t, err) {
			assert.Equal(t, b.Block, id)
		}
	}
}
