package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	DigitalGoodsListingType    = 3
	DigitalGoodsListingSubType = 0
)

type DigitalGoodsListing struct {
	*pb.DigitalGoodsListing
}

func (tx *DigitalGoodsListing) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint16(uint16(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes([]byte(tx.Attachment.Description))
	e.WriteUint16(uint16(len(tx.Attachment.Tags)))
	e.WriteBytes([]byte(tx.Attachment.Tags))
	e.WriteUint32(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *DigitalGoodsListing) AttachmentSizeInBytes() int {
	return 2 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description) + 2 + len(tx.Attachment.Tags) + 4 + 8
}

func (tx *DigitalGoodsListing) GetType() uint16 {
	return DigitalGoodsListingSubType<<8 | DigitalGoodsListingType
}
