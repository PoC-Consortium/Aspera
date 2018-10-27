package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type MultiSameOutCreation struct {
	*pb.MultiSameOutCreation
}

func (tx *MultiSameOutCreation) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint8(uint8(len(tx.Attachment.Recipients)))
	for _, recip := range tx.Attachment.Recipients {
		e.WriteUint64(recip)
	}

	return e.Bytes()
}

func (tx *MultiSameOutCreation) SizeInBytes() int {
	return HeaderSize(tx.Header) + 1 + len(tx.Attachment.Recipients)*8
}
