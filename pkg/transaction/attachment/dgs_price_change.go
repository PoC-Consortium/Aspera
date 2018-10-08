package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPriceChangeAttachment struct {
	Goods    uint64 `json:"goods,omitempty,string"`
	PriceNQT uint64 `json:"priceNQT,omitempty"`
	Version  int8   `struct:"-" json:"version.DigitalGoodsPriceChange,omitempty"`
}

func DgsPriceChangeAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment DgsPriceChangeAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8, err
}

func (attachment *DgsPriceChangeAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
