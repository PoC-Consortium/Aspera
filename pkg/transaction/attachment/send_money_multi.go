package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type Payment struct {
	Recip  uint64
	Amount uint64
}

type SendMoneyMultiTransaction struct {
	NumRecipsAndAmounts uint8 `struct:"uint8,sizeof=RecipsAndAmounts"`
	RecipsAndAmounts    []Payment
}

func SendMoneyMultiTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SendMoneyMultiTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 1 + len(tx.RecipsAndAmounts)*(8+8), err
}
