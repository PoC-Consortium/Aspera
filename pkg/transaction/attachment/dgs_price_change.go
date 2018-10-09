package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPriceChange struct {
	Goods    uint64 `json:"goods,omitempty,string"`
	PriceNQT uint64 `json:"priceNQT,omitempty"`
	Version  int8   `struct:"-" json:"version.DigitalGoodsPriceChange,omitempty"`
}

func (attachment *DgsPriceChange) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 8, err
}

func (attachment *DgsPriceChange) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
