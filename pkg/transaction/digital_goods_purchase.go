package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	DigitalGoodsPurchaseType    = 3
	DigitalGoodsPurchaseSubType = 4
)

type DigitalGoodsPurchase struct {
	*pb.DigitalGoodsPurchase
}

func EmptyDigitalGoodsPurchase() *DigitalGoodsPurchase {
	return &DigitalGoodsPurchase{
		DigitalGoodsPurchase: &pb.DigitalGoodsPurchase{
			Attachment: &pb.DigitalGoodsPurchase_Attachment{},
		},
	}
}

func (tx *DigitalGoodsPurchase) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint32(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
	e.WriteUint32(tx.Attachment.DeliveryDeadlineTimestamp)
}

func (tx *DigitalGoodsPurchase) AttachmentSizeInBytes() int {
	return 8 + 4 + 8 + 4
}

func (tx *DigitalGoodsPurchase) GetType() uint16 {
	return DigitalGoodsPurchaseSubType<<8 | DigitalGoodsPurchaseType
}

func (tx *DigitalGoodsPurchase) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *DigitalGoodsPurchase) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
