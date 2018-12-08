package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	MultiSameOutCreationType    = 0
	MultiSameOutCreationSubType = 2
)

type MultiSameOutCreation struct {
	*pb.MultiSameOutCreation
}

func EmptyMultiSameOutCreation() *MultiSameOutCreation {
	return &MultiSameOutCreation{
		MultiSameOutCreation: &pb.MultiSameOutCreation{
			Attachment: &pb.MultiSameOutCreation_Attachment{},
		},
	}
}

func (tx *MultiSameOutCreation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Recipients)))
	for _, recip := range tx.Attachment.Recipients {
		e.WriteUint64(recip)
	}
}

func (tx *MultiSameOutCreation) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Recipients = make([]uint64, d.ReadUint8())
	for i := range tx.Attachment.Recipients {
		tx.Attachment.Recipients[i] = d.ReadUint64()
	}
}

func (tx *MultiSameOutCreation) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Recipients)*8
}

func (tx *MultiSameOutCreation) GetType() uint16 {
	return MultiSameOutCreationSubType<<8 | MultiSameOutCreationType
}

func (tx *MultiSameOutCreation) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *MultiSameOutCreation) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
