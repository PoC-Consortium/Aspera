package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	DigitalGoodsListingType    = 3
	DigitalGoodsListingSubType = 0
)

type DigitalGoodsListing struct {
	*pb.DigitalGoodsListing
}

func EmptyDigitalGoodsListing() *DigitalGoodsListing {
	return &DigitalGoodsListing{
		DigitalGoodsListing: &pb.DigitalGoodsListing{
			Attachment: &pb.DigitalGoodsListing_Attachment{},
		},
	}
}

func (tx *DigitalGoodsListing) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint16(uint16(len(tx.Attachment.Name)))
	e.WriteBytes(tx.Attachment.Name)
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes(tx.Attachment.Description)
	e.WriteUint16(uint16(len(tx.Attachment.Tags)))
	e.WriteBytes(tx.Attachment.Tags)
	e.WriteUint32(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *DigitalGoodsListing) ReadAttachmentBytes(d encoding.Decoder) {
	nameLen := d.ReadUint16()
	tx.Attachment.Name = d.ReadBytes(int(nameLen))
	descriptionLen := d.ReadUint16()
	tx.Attachment.Description = d.ReadBytes(int(descriptionLen))
	tagsLen := d.ReadUint16()
	tx.Attachment.Tags = d.ReadBytes(int(tagsLen))
	tx.Attachment.Quantity = d.ReadUint32()
	tx.Attachment.Price = d.ReadUint64()
}

func (tx *DigitalGoodsListing) AttachmentSizeInBytes() int {
	return 2 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description) + 2 + len(tx.Attachment.Tags) + 4 + 8
}

func (tx *DigitalGoodsListing) GetType() uint16 {
	return DigitalGoodsListingSubType<<8 | DigitalGoodsListingType
}

func (tx *DigitalGoodsListing) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *DigitalGoodsListing) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
