package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	MultiOutCreationType    = 0
	MultiOutCreationSubType = 1
)

type MultiOutCreation struct {
	*pb.MultiOutCreation
}

func (tx *MultiOutCreation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Recipients)))
	for _, recipIdAndAmount := range tx.Attachment.Recipients {
		e.WriteUint64(recipIdAndAmount.Id)
		e.WriteUint64(recipIdAndAmount.Amount)
	}
}

func (tx *MultiOutCreation) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Recipients)*(8+8)
}

func (tx *MultiOutCreation) GetType() uint16 {
	return MultiOutCreationSubType<<8 | MultiOutCreationType
}
