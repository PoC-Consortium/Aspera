package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsQuantityChangeTransaction struct {
	Goods         uint64
	DeltaQuantity uint32
}

func DgsQuantityChangeTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx DgsQuantityChangeTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err

}
