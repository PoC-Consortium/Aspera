package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsRefundTransaction struct {
	Purchase  uint64
	RefundNQT uint64
}

func DgsRefundTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx DgsRefundTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8 + 8, err
}
