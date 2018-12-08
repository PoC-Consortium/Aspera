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

func EmptySubscriptionPayment() *SubscriptionPayment {
	return &SubscriptionPayment{
		SubscriptionPayment: &pb.SubscriptionPayment{
			Attachment: &pb.SubscriptionPayment_Attachment{},
		},
	}
}

func (tx *SubscriptionPayment) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
}

func (tx *SubscriptionPayment) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Id = d.ReadUint64()
}

func (tx *SubscriptionPayment) AttachmentSizeInBytes() int {
	return 8
}

func (tx *SubscriptionPayment) GetType() uint16 {
	return SubscriptionPaymentSubType<<8 | SubscriptionPaymentType
}

func (tx *SubscriptionPayment) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *SubscriptionPayment) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
