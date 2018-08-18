package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsRefundTransaction struct {
	Purchase uint64
}

func DgsRefundTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx DgsRefundTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
