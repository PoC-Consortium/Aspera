package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	BidOrderPlacementType    = 2
	BidOrderPlacementSubType = 3
)

type BidOrderPlacement struct {
	*pb.BidOrderPlacement
}

func (tx *BidOrderPlacement) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *BidOrderPlacement) AttachmentSizeInBytes() int {
	return 8 + 8 + 8
}

func (tx *BidOrderPlacement) GetType() uint16 {
	return BidOrderPlacementSubType<<8 | BidOrderPlacementType
}
