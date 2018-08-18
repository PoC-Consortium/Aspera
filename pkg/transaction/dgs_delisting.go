package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsDelistingTransaction struct {
	Goods uint64
}

func DgsDelistingTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx DgsDelistingTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, err
}
