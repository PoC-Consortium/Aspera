package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SellAliasTransaction struct {
	NumAlias uint8 `struct:"uint8,sizeof=Alias"`
	Alias    []byte
	PriceNQT int64
}

func SellAliasTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx SellAliasTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
