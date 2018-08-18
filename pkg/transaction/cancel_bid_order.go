package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type CancelBidOrderTransaction struct {
	Order uint64
}

func CancelBidOrderTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx CancelBidOrderTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
