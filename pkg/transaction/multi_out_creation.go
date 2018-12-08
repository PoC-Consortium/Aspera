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

func (tx *MultiOutCreation) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Recipients = make([]*pb.MultiOutCreation_Attachment_Recipients, d.ReadUint8())
	for i := range tx.Attachment.Recipients {
		tx.Attachment.Recipients[i] = &pb.MultiOutCreation_Attachment_Recipients{
			Id:     d.ReadUint64(),
			Amount: d.ReadUint64(),
		}
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
