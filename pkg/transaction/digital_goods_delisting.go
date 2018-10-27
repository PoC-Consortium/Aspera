package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsDelisting struct {
	*pb.DigitalGoodsDelisting
}

func (tx *DigitalGoodsDelisting) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Id)

	return e.Bytes()
}

func (tx *DigitalGoodsDelisting) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8
}
