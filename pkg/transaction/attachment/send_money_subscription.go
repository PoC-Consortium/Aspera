package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneySubscriptionTransaction struct {
	Frequency uint32
}

func SendMoneySubscriptionTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SendMoneySubscriptionTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 4, err
}
