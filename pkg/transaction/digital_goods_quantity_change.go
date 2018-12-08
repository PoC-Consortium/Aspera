package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	DigitalGoodsQuantityChangeType    = 3
	DigitalGoodsQuantityChangeSubType = 3
)

type DigitalGoodsQuantityChange struct {
	*pb.DigitalGoodsQuantityChange
}

func EmptyDigitalGoodsQuantityChange() *DigitalGoodsQuantityChange {
	return &DigitalGoodsQuantityChange{
		DigitalGoodsQuantityChange: &pb.DigitalGoodsQuantityChange{
			Attachment: &pb.DigitalGoodsQuantityChange_Attachment{},
		},
	}
}

func (tx *DigitalGoodsQuantityChange) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
	e.WriteInt32(tx.Attachment.Delta)
}

func (tx *DigitalGoodsQuantityChange) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Id = d.ReadUint64()
	tx.Attachment.Delta = d.ReadInt32()
}

func (tx *DigitalGoodsQuantityChange) AttachmentSizeInBytes() int {
	return 8 + 4
}

func (tx *DigitalGoodsQuantityChange) GetType() uint16 {
	return DigitalGoodsQuantityChangeSubType<<8 | DigitalGoodsQuantityChangeType
}

func (tx *DigitalGoodsQuantityChange) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *DigitalGoodsQuantityChange) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
