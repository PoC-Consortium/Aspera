package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type BidOrderPlacement struct {
	*pb.AskOrderPlacement
}

func (tx *BidOrderPlacement) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *BidOrderPlacement) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 8 + 8
}
