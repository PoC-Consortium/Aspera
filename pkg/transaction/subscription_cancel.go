package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	SubscriptionCancelType    = 21
	SubscriptionCancelSubType = 4
)

type SubscriptionCancel struct {
	*pb.SubscriptionCancel
}

func EmptySubscriptionCancel() *SubscriptionCancel {
	return &SubscriptionCancel{
		SubscriptionCancel: &pb.SubscriptionCancel{
			Attachment: &pb.SubscriptionCancel_Attachment{},
		},
	}
}

func (tx *SubscriptionCancel) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
}

func (tx *SubscriptionCancel) AttachmentSizeInBytes() int {
	return 8
}

func (tx *SubscriptionCancel) GetType() uint16 {
	return SubscriptionCancelSubType<<8 | SubscriptionCancelType
}
