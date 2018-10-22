package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsQuantityChange struct {
	*pb.DigitalGoodsQuantityChange
}

func (tx *DigitalGoodsQuantityChange) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.Goods)
	e.WriteInt32(tx.Attachment.DeltaQuantity)
}

func (tx *DigitalGoodsQuantityChange) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 4
}
