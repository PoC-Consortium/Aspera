package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type MultiOutCreation struct {
	*pb.MultiOutCreation
}

func (tx *MultiOutCreation) ToBytes() {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint8(uint8(len(tx.Attachment.Recipients)))
	for _, recipIdAndAmount := range tx.Attachment.Recipients {
		e.WriteUint64(recipAndId.Id)
		e.WriteUint64(recipAndId.Amount)
	}
}

func (tx *MultiOutCreation) SizeInBytes() {
	return HeaderSize(tx.TransactionHeader) + 1 + len(tx.Attachment.Recipients)*(8+8)
}
