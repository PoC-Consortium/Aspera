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

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint32(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
	e.WriteUint32(tx.Attachment.DeliveryDeadlineTimestamp)

	return e.Bytes()
}

func (tx *DigitalGoodsPurchase) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 4 + 8 + 4
}
