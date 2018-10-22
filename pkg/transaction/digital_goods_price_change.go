package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsPriceChange struct {
	*pb.DigitalGoodsPriceChange
}

func (tx *DigitalGoodsPriceChange) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.Goods)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *DigitalGoodsDelisting) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 8
}
