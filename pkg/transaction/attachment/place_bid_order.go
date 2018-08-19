package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceBidOrderTransaction struct {
	Asset       uint64
	QuantityQNT uint64
	PriceNQT    uint64
}

func PlaceBidOrderTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx PlaceBidOrderTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8 + 8 + 8, err
}
