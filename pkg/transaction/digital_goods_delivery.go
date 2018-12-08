package transaction

import (
	"encoding/hex"
	"math"

	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	DigitalGoodsDeliveryType    = 3
	DigitalGoodsDeliverySubType = 5
)

type DigitalGoodsDelivery struct {
	*pb.DigitalGoodsDelivery
}

func EmptyDigitalGoodsDelivery() *DigitalGoodsDelivery {
	return &DigitalGoodsDelivery{
		DigitalGoodsDelivery: &pb.DigitalGoodsDelivery{
			Attachment: &pb.DigitalGoodsDelivery_Attachment{},
		},
	}
}

func (tx *DigitalGoodsDelivery) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Purchase)
	// java wallet <3
	l := len(tx.Attachment.Data) / 2
	if tx.Attachment.IsText {
		e.WriteInt32(int32(l) | math.MinInt32)
	} else {
		e.WriteInt32(int32(l))
	}
	data := make([]byte, l)
	if _, err := hex.Decode(data, tx.Attachment.Data); err != nil {
		return
	}
	e.WriteBytes(data)
	e.WriteBytes(tx.Attachment.Nonce)
	e.WriteUint64(tx.Attachment.Discount)
}

func (tx *DigitalGoodsDelivery) AttachmentSizeInBytes() int {
	return 8 + 4 + len(tx.Attachment.Data)/2 + len(tx.Attachment.Nonce) + 8
}

func (tx *DigitalGoodsDelivery) GetType() uint16 {
	return DigitalGoodsDeliverySubType<<8 | DigitalGoodsDeliveryType
}
