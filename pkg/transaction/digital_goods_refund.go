package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	DigitalGoodsRefundType    = 3
	DigitalGoodsRefundSubType = 7
)

type DigitalGoodsRefund struct {
	*pb.DigitalGoodsRefund
}

func EmptyDigitalGoodsRefund() *DigitalGoodsRefund {
	return &DigitalGoodsRefund{
		DigitalGoodsRefund: &pb.DigitalGoodsRefund{
			Attachment: &pb.DigitalGoodsRefund_Attachment{},
		},
	}
}

func (tx *DigitalGoodsRefund) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Purchase)
	e.WriteUint64(tx.Attachment.Refund)
}

func (tx *DigitalGoodsRefund) AttachmentSizeInBytes() int {
	return 8 + 8
}

func (tx *DigitalGoodsRefund) GetType() uint16 {
	return DigitalGoodsRefundSubType<<8 | DigitalGoodsRefundType
}

func (tx *DigitalGoodsRefund) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *DigitalGoodsRefund) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
