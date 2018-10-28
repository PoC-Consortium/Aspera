package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	DigitalGoodsDelistingType    = 3
	DigitalGoodsDelistingSubType = 1
)

type DigitalGoodsDelisting struct {
	*pb.DigitalGoodsDelisting
}

func (tx *DigitalGoodsDelisting) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
}

func (tx *DigitalGoodsDelisting) AttachmentSizeInBytes() int {
	return 8
}

func (tx *DigitalGoodsDelisting) GetType() uint16 {
	return DigitalGoodsDelistingSubType<<8 | DigitalGoodsDelistingType
}
