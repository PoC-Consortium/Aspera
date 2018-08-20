package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceAskOrderAttachment struct {
	Asset       uint64
	QuantityQNT uint64
	PriceNQT    uint64
}

func PlaceAskOrderAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment PlaceAskOrderAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8 + 8, err
}

func (attachment *PlaceAskOrderAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
