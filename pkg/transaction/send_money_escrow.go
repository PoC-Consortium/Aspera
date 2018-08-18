package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyEscrowTransaction struct {
	AmountNQT      uint64
	EscrowDeadline uint32
	DeadlineAction uint8
	NumSignees     uint8 `struct:"uint8,sizeof=Signees"`
	TotalSignees   uint8
	Signees        []uint64
}

func SendMoneyEscrowTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx SendMoneyEscrowTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
