package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceBidOrderTransaction struct {
	Asset       uint64
	QuantityQNT uint64
	PriceNQT    uint64
}

func PlaceBidOrderTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx PlaceBidOrderTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
