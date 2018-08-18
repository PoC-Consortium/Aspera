package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type BuyAliasTransaction struct {
	NumAlias uint8 `struct:"uint8,sizeof=Alias"`
	Alias    []byte
}

func BuyAliasTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx BuyAliasTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
