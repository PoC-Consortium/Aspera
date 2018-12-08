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

func (tx *SubscriptionCancel) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Id = d.ReadUint64()
}

func (tx *SubscriptionCancel) GetType() uint16 {
	return SubscriptionCancelSubType<<8 | SubscriptionCancelType
}

func (tx *SubscriptionCancel) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *SubscriptionCancel) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
