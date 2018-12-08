package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	DigitalGoodsPriceChangeType    = 3
	DigitalGoodsPriceChangeSubType = 2
)

type DigitalGoodsPriceChange struct {
	*pb.DigitalGoodsPriceChange
}

func EmptyDigitalGoodsPriceChange() *DigitalGoodsPriceChange {
	return &DigitalGoodsPriceChange{
		DigitalGoodsPriceChange: &pb.DigitalGoodsPriceChange{
			Attachment: &pb.DigitalGoodsPriceChange_Attachment{},
		},
	}
}

func (tx *DigitalGoodsPriceChange) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *DigitalGoodsPriceChange) AttachmentSizeInBytes() int {
	return 8 + 8
}

func (tx *DigitalGoodsPriceChange) GetType() uint16 {
	return DigitalGoodsPriceChangeSubType<<8 | DigitalGoodsPriceChangeType
}
