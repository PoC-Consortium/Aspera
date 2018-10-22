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

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.Purchase)
	e.WriteBytesWithInt32Len(tx.Attachment.IsText, tx.Attachment.GoodsData)
	e.WriteBytes(tx.Attachment.IsText, tx.Attachment.GoodsNonce)
	e.WriteUint64(tx.Attachment.Discount)
}

func (tx *DigitalGoodsDelivery) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 4 + len(tx.Attachment.GoodsData) +
		len(tx.Attachment.GoodsNonce) + 8
}
