package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type TransferAssetTransaction struct {
	Asset       uint64
	QuantityQNT uint64
}

func TransferAssetTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx TransferAssetTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8 + 8, err
}
