package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type TransferAssetTransaction struct {
	Asset       uint64
	QuantityQNT uint64
}

func TransferAssetTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx TransferAssetTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
