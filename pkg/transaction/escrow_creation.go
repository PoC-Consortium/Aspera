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

func (tx *EscrowCreation) AttachmentSizeInBytes() int {
	return 8 + 4 + 1 + 1 + 1 + len(tx.Attachment.Signers)*8
}

func (tx *EscrowCreation) GetType() uint16 {
	return EscrowCreationSubType<<8 | EscrowCreationType
}
