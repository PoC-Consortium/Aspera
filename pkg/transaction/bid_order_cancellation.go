package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	BidOrderCancellationType    = 2
	BidOrderCancellationSubType = 5
)

type BidOrderCancellation struct {
	*pb.BidOrderCancellation
}

func (tx *BidOrderCancellation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Order)
}

func (tx *BidOrderCancellation) AttachmentSizeInBytes() int {
	return 8
}

func (tx *BidOrderCancellation) GetType() uint16 {
	return BidOrderCancellationSubType<<8 | BidOrderCancellationType
}
