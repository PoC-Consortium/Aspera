package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsPurchase struct {
	*pb.DigitalGoodsPurchase
}

func (tx *DigitalGoodsPurchase) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.Goods)
	e.WriteUint32(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
	e.WriteUint32(tx.Attachment.DeliveryDeadlineTimestamp)
}

func (tx *DigitalGoodsPurchase) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 4 + 8 + 4
}
