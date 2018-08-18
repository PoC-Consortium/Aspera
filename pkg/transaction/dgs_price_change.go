package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPriceChangeTransaction struct {
	Goods    uint64
	PriceNQT uint64
}

func DgsPriceChangeTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx DgsPriceChangeTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
