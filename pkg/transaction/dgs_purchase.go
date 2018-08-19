package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPurchaseTransaction struct {
	Goods                     uint64
	Quantity                  uint32
	PriceNQT                  uint64
	DeliveryDeadlineTimestamp uint32
}

func DgsPurchaseTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx DgsPurchaseTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8 + 4 + 8 + 4, err
}
