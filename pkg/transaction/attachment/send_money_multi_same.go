package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyMultiSameTransaction struct {
	RecipCount uint8 `struct:"uint8,sizeof=Recips"`
	Recips     []uint64
}

func SendMoneyMultiSameTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SendMoneyMultiSameTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 1 + len(tx.Recips)*8, err
}
