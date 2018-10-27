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

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)

	return e.Bytes()
}

func (tx *AskOrderPlacement) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 8 + 8
}
