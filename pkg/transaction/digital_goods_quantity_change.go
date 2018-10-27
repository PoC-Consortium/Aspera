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

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Id)
	e.WriteInt32(tx.Attachment.Delta)

	return e.Bytes()
}

func (tx *DigitalGoodsQuantityChange) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 4
}
