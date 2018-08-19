package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type LeaseBalanceTransaction struct {
	Period uint16
}

func LeaseBalanceTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx LeaseBalanceTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 2, err
}
