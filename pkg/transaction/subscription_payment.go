package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	SubscriptionPaymentType    = 21
	SubscriptionPaymentSubType = 5
)

type SubscriptionPayment struct {
	*pb.SubscriptionPayment
}

func (tx *SubscriptionPayment) WriteAttachmentBytes(e encoding.Encoder) {
	for _, id := range tx.Attachment.Ids {
		e.WriteUint64(id)
	}
}

func (tx *SubscriptionPayment) AttachmentSizeInBytes() int {
	return len(tx.Attachment.Ids) * 8
}

func (tx *SubscriptionPayment) GetType() uint16 {
	return SubscriptionPaymentSubType<<8 | SubscriptionPaymentType
}
