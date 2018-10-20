package block

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ac0v/aspera/pkg/transaction"
)

func TestGetID(t *testing.T) {
	for _, blockTest := range BlockTests {
		b := blockTest.Block

		// dummy set, because we need to hash over number of transactions
		b.Transactions = make([]*transaction.Transaction, blockTest.TXLen)
		if blockTest.BlockATs != nil {
			b.BlockATs = &blockTest.BlockATs
		}

		_, id, err := b.CalculateHashAndID()
		if assert.Nil(t, err) {
			assert.Equal(t, b.Block, id)
		}
	}
}

func TestCalculateGenerationSignature(t *testing.T) {
	genSig, _ := hex.DecodeString("c26ef60f51aa5fc6225a481f08e51903085067a8a7d558f94712d702f2a67bb4")
	genPubKey, _ := hex.DecodeString("735f5b8e04f45080acdcb2b7ecfc19697c4022f13f5e713afeed537710dd1529")
	b := &Block{
		GenerationSignature: genSig,
		GeneratorPublicKey:  genPubKey,
	}

	nextGenSig := CalculateGenerationSignature(b)
	assert.Equal(t, "a62b500a5dfc7f5e614fcf4917d83ffccd01e9c9643d7d0e982c75043d27baff",
		hex.EncodeToString(nextGenSig))
}
