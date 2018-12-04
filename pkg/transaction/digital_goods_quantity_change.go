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

func (tx *DigitalGoodsQuantityChange) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
	e.WriteInt32(tx.Attachment.Delta)
}

func (tx *DigitalGoodsQuantityChange) AttachmentSizeInBytes() int {
	return 8 + 4
}

func (tx *DigitalGoodsQuantityChange) GetType() uint16 {
	return DigitalGoodsQuantityChangeSubType<<8 | DigitalGoodsQuantityChangeType
}
