package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsDelivery struct {
	*pb.DigitalGoodsDelivery
}

func (tx *DigitalGoodsDelivery) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Purchase)
	e.WriteBytesWithInt32Len(tx.Attachment.IsText, tx.Attachment.Data)
	e.WriteBytes(tx.Attachment.Nonce)
	e.WriteUint64(tx.Attachment.Discount)

	return e.Bytes()
}

func (tx *DigitalGoodsDelivery) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 4 + len(tx.Attachment.Data) + len(tx.Attachment.Nonce) + 8
}
