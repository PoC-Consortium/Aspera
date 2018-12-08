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

func EmptySubscriptionSubscribe() *SubscriptionSubscribe {
	return &SubscriptionSubscribe{
		SubscriptionSubscribe: &pb.SubscriptionSubscribe{
			Attachment: &pb.SubscriptionSubscribe_Attachment{},
		},
	}
}

func (tx *SubscriptionSubscribe) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint32(tx.Attachment.Frequency)
}
func (tx *SubscriptionSubscribe) AttachmentSizeInBytes() int {
	return 4
}

func (tx *SubscriptionSubscribe) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Frequency = d.ReadUint32()
}

func (tx *SubscriptionSubscribe) GetType() uint16 {
	return SubscriptionSubscribeSubType<<8 | SubscriptionSubscribeType
}

func (tx *SubscriptionSubscribe) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *SubscriptionSubscribe) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
