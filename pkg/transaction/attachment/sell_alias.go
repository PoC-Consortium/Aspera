package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SellAlias struct {
	NumAlias uint8  `struct:"uint8,sizeof=Alias" json:"-"`
	Alias    string `json:"alias"`
	PriceNQT int64  `json:"priceNQT"`
	Version  int8   `struct:"-" json:"version.AliasSell,omitempty"`
}

func (attachment *SellAlias) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 1 + len(attachment.Alias) + 8, err
}

func (attachment *SellAlias) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
