package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type BuyAlias struct {
	NumAlias uint8  `struct:"uint8,sizeof=Alias" json:"-"`
	Alias    string `json:"alias"`
	Version  int8   `struct:"-" json:"version.AliasAssignment,omitempty"`
}

func (attachment *BuyAlias) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 1 + len(attachment.Alias), err
}

func (attachment *BuyAlias) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
