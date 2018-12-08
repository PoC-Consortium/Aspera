package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	EscrowSignType    = 21
	EscrowSignSubType = 1
)

type EscrowSign struct {
	*pb.EscrowSign
}

func EmptyEscrowSign() *EscrowSign {
	return &EscrowSign{
		EscrowSign: &pb.EscrowSign{
			Attachment: &pb.EscrowSign_Attachment{},
		},
	}
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
