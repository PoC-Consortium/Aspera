package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	DigitalGoodsFeedbackType    = 3
	DigitalGoodsFeedbackSubType = 6
)

type DigitalGoodsFeedback struct {
	*pb.DigitalGoodsFeedback
}

func (tx *DigitalGoodsFeedback) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Purchase)
}

func (tx *DigitalGoodsFeedback) AttachmentSizeInBytes() int {
	return 8
}

func (tx *DigitalGoodsFeedback) GetType() uint16 {
	return DigitalGoodsFeedbackSubType<<8 | DigitalGoodsFeedbackType
}
