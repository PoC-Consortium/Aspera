package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type BidOrderCancellation struct {
	*pb.BidOrderCancellation
}

func (tx *BidOrderCancellation) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Order)

	return e.Bytes()
}

func (tx *BidOrderCancellation) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8
}
