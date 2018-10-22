package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type DigitalGoodsListing struct {
	*pb.DigitalGoodsListing
}

func (tx *DigitalGoodsListing) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint16(uint16(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes([]byte(tx.Attachment.Description))
	e.WriteUint16(uint16(len(tx.Attachment.Tags)))
	e.WriteBytes([]byte(tx.Attachment.Tags))
	e.WriteUint32(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *DigitalGoodsListing) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 2 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description) +
		2 + len(tx.Attachment.Tags) + 4 + 8
}
