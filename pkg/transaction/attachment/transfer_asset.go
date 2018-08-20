package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type TransferAssetAttachment struct {
	Asset       uint64
	QuantityQNT uint64
}

func TransferAssetAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment TransferAssetAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8, err
}

func (attachment *TransferAssetAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
