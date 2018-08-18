package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type EscrowResultTransaction struct {
	EscrowID uint64
	Decision uint8
}

func EscrowResultTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx EscrowResultTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
