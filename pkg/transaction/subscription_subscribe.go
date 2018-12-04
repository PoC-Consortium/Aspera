package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	SubscriptionSubscribeType    = 21
	SubscriptionSubscribeSubType = 3
)

type SubscriptionSubscribe struct {
	*pb.SubscriptionSubscribe
}

func (tx *SubscriptionSubscribe) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint32(tx.Attachment.Frequency)
}
func (tx *SubscriptionSubscribe) AttachmentSizeInBytes() int {
	return 4
}

func (tx *SubscriptionSubscribe) GetType() uint16 {
	return SubscriptionSubscribeSubType<<8 | SubscriptionSubscribeType
}
