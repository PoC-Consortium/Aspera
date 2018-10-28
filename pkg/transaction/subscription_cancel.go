package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	SubscriptionCancelType    = 21
	SubscriptionCancelSubType = 4
)

type SubscriptionCancel struct {
	*pb.SubscriptionCancel
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
