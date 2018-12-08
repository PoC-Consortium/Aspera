package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	EscrowCreationType    = 21
	EscrowCreationSubType = 0
)

type EscrowCreation struct {
	*pb.EscrowCreation
}

func EmptyEscrowCreation() *EscrowCreation {
	return &EscrowCreation{
		EscrowCreation: &pb.EscrowCreation{
			Attachment: &pb.EscrowCreation_Attachment{},
		},
	}
}

func (tx *EscrowCreation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Amount)
	e.WriteUint32(tx.Attachment.Deadline)
	e.WriteUint8(uint8(tx.Attachment.DeadlineAction))
	e.WriteUint8(uint8(tx.Attachment.RequiredSigners))
	e.WriteUint8(uint8(len(tx.Attachment.Signers)))
	for _, signer := range tx.Attachment.Signers {
		e.WriteUint64(signer)
	}
}

func (tx *EscrowCreation) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Amount = d.ReadUint64()
	tx.Attachment.Deadline = d.ReadUint32()
	tx.Attachment.DeadlineAction = pb.DeadlineAction((d.ReadUint8()))
	tx.Attachment.RequiredSigners = int32(d.ReadUint8())
	tx.Attachment.Signers = make([]uint64, d.ReadUint8())
	for i := range tx.Attachment.Signers {
		tx.Attachment.Signers[i] = d.ReadUint64()
	}
}

func (tx *EscrowCreation) AttachmentSizeInBytes() int {
	return 8 + 4 + 1 + 1 + 1 + len(tx.Attachment.Signers)*8
}

func (tx *EscrowCreation) GetType() uint16 {
	return EscrowCreationSubType<<8 | EscrowCreationType
}

func (tx *EscrowCreation) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *EscrowCreation) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
