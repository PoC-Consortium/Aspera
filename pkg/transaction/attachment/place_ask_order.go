package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceAskOrder struct {
	Asset       uint64 `json:"asset,omitempty,string"`
	QuantityQNT uint64 `json:"quantityQNT,omitempty"`
	PriceNQT    uint64 `json:"priceNQT,omitempty"`

	Comment string `struct:"-" json:"comment,omitempty"`
	Version int8   `struct:"-" json:"version.AskOrderPlacement,omitempty"`
}

func (attachment *PlaceAskOrder) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 8 + 8, err
}

func (attachment *PlaceAskOrder) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *PlaceAskOrder) GetFlag() uint32 {
	return StandardAttachmentFlag
}
