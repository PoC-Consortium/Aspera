package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyMultiSame struct {
	RecipCount uint8             `struct:"uint8,sizeof=Recips" json:"-"`
	Recips     UInt64StringSlice `json:"recipients"` // fix restruct
	Version    int8              `struct:"-" json:"version.MultiSameOutCreation,omitempty"`
}

func (attachment *SendMoneyMultiSame) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 1 + len(attachment.Recips)*8, err
}

func (attachment *SendMoneyMultiSame) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
