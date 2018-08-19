package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SubscriptionCancelTransaction struct {
	Subscription uint64
}

func SubscriptionCancelTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SubscriptionCancelTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8, err
}
