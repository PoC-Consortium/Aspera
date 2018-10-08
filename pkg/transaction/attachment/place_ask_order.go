package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceAskOrderAttachment struct {
	Asset       uint64 `json:"asset,omitempty,string"`
	QuantityQNT uint64 `json:"quantityQNT,omitempty"`
	PriceNQT    uint64 `json:"priceNQT,omitempty"`

	Comment string `struct:"-" json:"comment,omitempty"`
	Version int8   `struct:"-" json:"version.AskOrderPlacement,omitempty"`
}

func PlaceAskOrderAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment PlaceAskOrderAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8 + 8, err
}

func (attachment *PlaceAskOrderAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
