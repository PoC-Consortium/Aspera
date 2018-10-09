package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type Payment struct {
	Recip  uint64 `json:",string"`
	Amount uint64 `json:",string"`
}

type SendMoneyMulti struct {
	NumRecipsAndAmounts uint8     `struct:"uint8,sizeof=RecipsAndAmounts" json:"-"`
	RecipsAndAmounts    []Payment `json:"recipients"` // fix restruct
	Version             int8      `struct:"-" json:"version.MultiOutCreation,omitempty"`
}

func (attachment *SendMoneyMulti) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 1 + len(attachment.RecipsAndAmounts)*(8+8), err
}

func (attachment *SendMoneyMulti) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
