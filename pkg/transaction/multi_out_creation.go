package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	MultiOutCreationType    = 0
	MultiOutCreationSubType = 1
)

type MultiOutCreation struct {
	*pb.MultiOutCreation
}

func EmptyMultiOutCreation() *MultiOutCreation {
	return &MultiOutCreation{
		MultiOutCreation: &pb.MultiOutCreation{
			Attachment: &pb.MultiOutCreation_Attachment{},
		},
	}
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

func (tx *MultiOutCreation) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *MultiOutCreation) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
