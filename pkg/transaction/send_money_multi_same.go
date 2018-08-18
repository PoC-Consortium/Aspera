package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyMultiSameTransaction struct {
	RecipCount uint8 `struct:"uint8,sizeof=Recips"`
	Recips     []uint64
}

func SendMoneyMultiSameTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx SendMoneyMultiSameTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
