package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type LeaseBalance struct {
	Period  uint16 `json:"period"`
	Version int8   `struct:"-" json:"version.EffectiveBalanceLeasing,omitempty"`
}

func (attachment *LeaseBalance) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 2, err
}

func (attachment *LeaseBalance) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
