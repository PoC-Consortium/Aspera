package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsQuantityChangeAttachment struct {
	Goods         uint64 `json:"goods,omitempty,string"`
	DeltaQuantity uint32 `json:"deltaQuantity,omitempty"`
	Version       int8   `struct:"-" json:"version.DigitalGoodsQuantityChange,omitempty"`
}

func DgsQuantityChangeAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment DgsQuantityChangeAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 4, err
}

func (attachment *DgsQuantityChangeAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
