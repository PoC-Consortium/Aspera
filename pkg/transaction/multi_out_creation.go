package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type MultiOutCreation struct {
	*pb.MultiOutCreation
}

func (tx *MultiOutCreation) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint8(uint8(len(tx.Attachment.Recipients)))
	for _, recipIdAndAmount := range tx.Attachment.Recipients {
		e.WriteUint64(recipIdAndAmount.Id)
		e.WriteUint64(recipIdAndAmount.Amount)
	}

	return e.Bytes()
}

func (tx *MultiOutCreation) SizeInBytes() int {
	return HeaderSize(tx.Header) + 1 + len(tx.Attachment.Recipients)*(8+8)
}
