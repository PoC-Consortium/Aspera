package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	AssetTransferType    = 2
	AssetTransferSubType = 2
)

type AssetTransfer struct {
	*pb.AssetTransfer
}

func (tx *AssetTransfer) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
}

func (tx *AssetTransfer) AttachmentSizeInBytes() int {
	return 8 + 8
}

func (tx *AssetTransfer) GetType() uint16 {
	return AssetTransferSubType<<8 | AssetTransferType
}
