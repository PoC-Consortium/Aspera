package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsQuantityChangeTransaction struct {
	Goods         uint64
	DeltaQuantity uint32
}

func DgsQuantityChangeTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx DgsQuantityChangeTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8 + 4, err
}
