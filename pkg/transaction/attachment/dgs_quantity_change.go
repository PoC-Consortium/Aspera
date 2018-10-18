package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsQuantityChange struct {
	Goods         uint64 `json:"goods,omitempty,string"`
	DeltaQuantity uint32 `json:"deltaQuantity,omitempty"`
	Version       int8   `struct:"-" json:"version.DigitalGoodsQuantityChange,omitempty"`
}

func (attachment *DgsQuantityChange) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 4, err
}

func (attachment *DgsQuantityChange) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *DgsQuantityChange) GetFlag() uint32 {
	return StandardAttachmentFlag
}
