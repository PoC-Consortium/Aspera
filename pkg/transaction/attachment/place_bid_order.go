package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PlaceBidOrder struct {
	Asset       uint64 `json:"asset,omitempty,string"`
	QuantityQNT uint64 `json:"quantityQNT,omitempty"`
	PriceNQT    uint64 `json:"priceNQT,omitempty"`
	Version     int8   `struct:"-" json:"version.BidOrderPlacement,omitempty"`
}

func (attachment *PlaceBidOrder) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 8 + 8, err
}

func (attachment *PlaceBidOrder) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *PlaceBidOrder) GetFlag() uint32 {
	return StandardAttachmentFlag
}
