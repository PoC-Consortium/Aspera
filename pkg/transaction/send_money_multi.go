package transaction

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

func SendMoneyMultiTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx SendMoneyMultiTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
