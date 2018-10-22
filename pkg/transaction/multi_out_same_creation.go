package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type MultiOutSameCreation struct {
	*pb.MultiOutSameCreation
}

func (tx *MultiOutSameCreation) ToBytes() {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint8(uint8(len(tx.Attachment.Recipients)))
	for _, recip := range tx.Attachment.Recipients {
		e.WriteUint64(recip)
	}
}

func (tx *MultiSameOutCreation) SizeInBytes() {
	return HeaderSize(tx.TransactionHeader) + 1 + len(tx.Attachment.Recipients)*8
}
