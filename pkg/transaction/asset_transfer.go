package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	AssetTransferType    = 2
	AssetTransferSubType = 1
)

type AssetTransfer struct {
	*pb.AssetTransfer
}

func (tx *AssetTransfer) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
	if tx.Header.Version == 0 {
		e.WriteUint16(uint16(len(tx.Attachment.Comment)))
		e.WriteBytes(tx.Attachment.Comment)
	}
}

func (tx *AssetTransfer) AttachmentSizeInBytes() int {
	if tx.Header.Version == 0 {
		return 8 + 8 + 2 + len(tx.Attachment.Comment)
	} else {
		return 8 + 8
	}
}

func (tx *AssetTransfer) GetType() uint16 {
	return AssetTransferSubType<<8 | AssetTransferType
}
