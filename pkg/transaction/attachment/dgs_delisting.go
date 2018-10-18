package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsDelisting struct {
	Goods   uint64 `json:"goods,string"`
	Version int8   `struct:"-" json:"version.DigitalGoodsDelisting,omitempty"`
}

func (attachment *DgsDelisting) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8, err
}

func (attachment *DgsDelisting) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *DgsDelisting) GetFlag() uint32 {
	return StandardAttachmentFlag
}
