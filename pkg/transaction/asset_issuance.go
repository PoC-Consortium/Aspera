package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	AssetIssuanceType    = 2
	AssetIssuanceSubType = 0
)

type AssetIssuance struct {
	*pb.AssetIssuance
}

func (tx *AssetIssuance) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes(tx.Attachment.Name)
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes(tx.Attachment.Description)
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint8(uint8(tx.Attachment.Decimals))
}

func (tx *AssetIssuance) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description) + 8 + 1
}

func (tx *AssetIssuance) GetType() uint16 {
	return AssetIssuanceSubType<<8 | AssetIssuanceType
}
