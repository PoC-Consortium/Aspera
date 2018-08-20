package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceBidOrderAttachment struct {
	Asset       uint64
	QuantityQNT uint64
	PriceNQT    uint64
}

func PlaceBidOrderAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment PlaceBidOrderAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8 + 8, err
}

func (attachment *PlaceBidOrderAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
