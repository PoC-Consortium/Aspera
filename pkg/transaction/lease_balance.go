package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type LeaseBalanceTransaction struct {
	Period uint16
}

func LeaseBalanceTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx LeaseBalanceTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
