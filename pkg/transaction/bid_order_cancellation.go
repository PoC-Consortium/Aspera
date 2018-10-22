package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type BidOrderCancellation struct {
	*pb.BidOrderCancellation
}

func (tx *BidOrderPlacement) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.Order)
}

func (tx *BidOrderCancellation) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8
}
