package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceAskOrderTransaction struct {
	Asset       uint64
	QuantityQNT uint64
	PriceNQT    uint64
}

func PlaceAskOrderTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx PlaceAskOrderTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
