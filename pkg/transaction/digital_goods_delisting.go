package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	DigitalGoodsDelistingType    = 3
	DigitalGoodsDelistingSubType = 1
)

type DigitalGoodsDelisting struct {
	*pb.DigitalGoodsDelisting
}

func EmptyDigitalGoodsDelisting() *DigitalGoodsDelisting {
	return &DigitalGoodsDelisting{
		DigitalGoodsDelisting: &pb.DigitalGoodsDelisting{
			Attachment: &pb.DigitalGoodsDelisting_Attachment{},
		},
	}
}

func (tx *DigitalGoodsDelisting) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
}

func (tx *DigitalGoodsDelisting) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Id = d.ReadUint64()
}

func (tx *DigitalGoodsDelisting) AttachmentSizeInBytes() int {
	return 8
}

func (tx *DigitalGoodsDelisting) GetType() uint16 {
	return DigitalGoodsDelistingSubType<<8 | DigitalGoodsDelistingType
}

func (tx *DigitalGoodsDelisting) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *DigitalGoodsDelisting) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
