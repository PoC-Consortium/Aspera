package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsFeedback struct {
	*pb.DigitalGoodsFeedback
}

func (tx *DigitalGoodsFeedback) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Purchase)

	return e.Bytes()
}

func (tx *DigitalGoodsFeedback) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8
}
