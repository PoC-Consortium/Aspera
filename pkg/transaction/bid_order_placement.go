package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	BidOrderPlacementType    = 2
	BidOrderPlacementSubType = 3
)

type BidOrderPlacement struct {
	*pb.BidOrderPlacement
}

func EmptyBidOrderPlacement() *BidOrderPlacement {
	return &BidOrderPlacement{
		BidOrderPlacement: &pb.BidOrderPlacement{
			Attachment: &pb.BidOrderPlacement_Attachment{},
		},
	}
}

func (tx *BidOrderPlacement) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *BidOrderPlacement) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Asset = d.ReadUint64()
	tx.Attachment.Quantity = d.ReadUint64()
	tx.Attachment.Price = d.ReadUint64()
}

func (tx *BidOrderPlacement) AttachmentSizeInBytes() int {
	return 8 + 8 + 8
}

func (tx *BidOrderPlacement) GetType() uint16 {
	return BidOrderPlacementSubType<<8 | BidOrderPlacementType
}

func (tx *BidOrderPlacement) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *BidOrderPlacement) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
