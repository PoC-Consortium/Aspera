package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type EscrowSignTransaction struct {
	Escrow   uint64
	Decision uint8
}

func EscrowSignTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx EscrowSignTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8 + 1, err
}
