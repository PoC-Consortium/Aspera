package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	DigitalGoodsPriceChangeType    = 3
	DigitalGoodsPriceChangeSubType = 2
)

type DigitalGoodsPriceChange struct {
	*pb.DigitalGoodsPriceChange
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
