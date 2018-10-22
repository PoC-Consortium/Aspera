package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsRefund struct {
	*pb.DigitalGoodsRefund
}

func (tx *DigitalGoodsRefund) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.Purchase)
	e.WriteUint64(tx.Attachment.RefundPrice)
}

func (tx *DigitalGoodsRefund) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 8
}
