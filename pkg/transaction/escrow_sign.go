package transaction

import (
	pb "github.com/PoC-Consortium/aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/aspera/pkg/encoding"
)

const (
	EscrowSignType    = 21
	EscrowSignSubType = 1
)

type EscrowSign struct {
	*pb.EscrowSign
}

func (tx *EscrowSign) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint8(uint8(tx.Attachment.Decision))
}

func (tx *EscrowSign) AttachmentSizeInBytes() int {
	return 8 + 1
}

func (tx *EscrowSign) GetType() uint16 {
	return EscrowSignSubType<<8 | EscrowSignType
}
