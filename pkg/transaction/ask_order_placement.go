package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AskOrderPlacement struct {
	*pb.AskOrderPlacement
}

func (tx *AskOrderPlacement) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.AssetTransfer)
	e.WriteUint64(tx.Attachment.Qantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *AskOrderPlacement) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 8 + 8
}
