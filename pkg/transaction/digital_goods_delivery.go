package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	DigitalGoodsDeliveryType    = 3
	DigitalGoodsDeliverySubType = 5
)

type DigitalGoodsDelivery struct {
	*pb.DigitalGoodsDelivery
}

func (tx *DigitalGoodsDelivery) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Purchase)
	e.WriteBytesWithInt32Len(tx.Attachment.IsText, tx.Attachment.Data)
	e.WriteBytes(tx.Attachment.Nonce)
	e.WriteUint64(tx.Attachment.Discount)
}

func (tx *DigitalGoodsDelivery) AttachmentSizeInBytes() int {
	return 8 + 4 + len(tx.Attachment.Data) + len(tx.Attachment.Nonce) + 8
}

func (tx *DigitalGoodsDelivery) GetType() uint16 {
	return DigitalGoodsDeliverySubType<<8 | DigitalGoodsDeliveryType
}
