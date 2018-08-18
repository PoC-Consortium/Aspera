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

func DgsPurchaseTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx DgsPurchaseTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
