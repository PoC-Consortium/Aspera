package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsFeedbackTransaction struct {
	Purchase uint64
}

func DgsFeedbackTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx DgsFeedbackTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 8, err
}
