package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type CancelAskOrderTransaction struct {
	Order uint64
}

func CancelAskOrderTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx CancelAskOrderTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8, err
}
