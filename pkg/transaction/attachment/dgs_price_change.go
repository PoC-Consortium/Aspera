package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPriceChangeTransaction struct {
	Goods    uint64
	PriceNQT uint64
}

func DgsPriceChangeTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx DgsPriceChangeTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8 + 8, err
}
