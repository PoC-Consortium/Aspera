package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	AskOrderCancellationType    = 2
	AskOrderCancellationSubType = 4
)

type AskOrderCancellation struct {
	*pb.AskOrderCancellation
}

func EmptyAskOrderCancellation() *AskOrderCancellation {
	return &AskOrderCancellation{
		AskOrderCancellation: &pb.AskOrderCancellation{
			Attachment: &pb.AskOrderCancellation_Attachment{},
		},
	}
}

func (tx *AskOrderCancellation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Order)
}

func (tx *AskOrderCancellation) AttachmentSizeInBytes() int {
	return 8
}

func (tx *AskOrderCancellation) GetType() uint16 {
	return AskOrderCancellationSubType<<8 | AskOrderCancellationType
}
