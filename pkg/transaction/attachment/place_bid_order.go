package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceBidOrderAttachment struct {
	Asset       uint64 `json:"asset,omitempty,string"`
	QuantityQNT uint64 `json:"quantityQNT,omitempty"`
	PriceNQT    uint64 `json:"priceNQT,omitempty"`
	Version     int8   `struct:"-" json:"version.BidOrderPlacement,omitempty"`
}

func PlaceBidOrderAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment PlaceBidOrderAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8 + 8, err
}

func (attachment *PlaceBidOrderAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
