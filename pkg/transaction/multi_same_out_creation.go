package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	MultiSameOutCreationType    = 0
	MultiSameOutCreationSubType = 2
)

type MultiSameOutCreation struct {
	*pb.MultiSameOutCreation
}

func (tx *MultiSameOutCreation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Recipients)))
	for _, recip := range tx.Attachment.Recipients {
		e.WriteUint64(recip)
	}
}

func (tx *MultiSameOutCreation) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Recipients)*8
}

func (tx *MultiSameOutCreation) GetType() uint16 {
	return MultiSameOutCreationSubType<<8 | MultiSameOutCreationType
}
