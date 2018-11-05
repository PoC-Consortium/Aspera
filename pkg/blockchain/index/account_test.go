package index

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/ac0v/aspera/pkg/account"

	"github.com/stretchr/testify/assert"
)

var ai *accountIndex

func TestMain(m *testing.M) {
	if err := os.RemoveAll("var"); err != nil {
		panic(err)
	}
	if err := os.Mkdir("var", 0777); err != nil {
		panic(err)
	}
	ai = NewAccountIndex().(*accountIndex)
	m.Run()
}

func TestIndex(t *testing.T) {
	a := account.NewAccount(7900104405094198526)
	a.PublicKey, _ = hex.DecodeString("d37ebb299c54fa4603f1b656f6bcf70810a9cd1d56e2ea979c7933d94ae3602a")

	err := ai.Index(a)
	if err != nil {
		panic(err)
	}

	expectedId := []byte(fmt.Sprintf("%020d", a.Id))
	id, err := ai.ByAddress("8V9Y-58B4-RVWP-8HQAV")
	if assert.Nil(t, err) {
		assert.Equal(t, expectedId, id)
	}
	id, err = ai.ByPublicKey(a.PublicKey)
	if assert.Nil(t, err) {
		assert.Equal(t, expectedId, id)
	}
	id = ai.ById(a.Id)
	assert.Equal(t, expectedId, id)
}
