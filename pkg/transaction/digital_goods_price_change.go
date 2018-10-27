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

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint64(tx.Attachment.Price)

	return e.Bytes()
}

func (tx *DigitalGoodsPriceChange) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 8
}
