package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type CancelAskOrder struct {
	Order   uint64 `json:"order,omitempty,string"`
	Version int8   `struct:"-" json:"version.AskOrderCancellation,omitempty"`
}

func (attachment *CancelAskOrder) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8, err
}

func (attachment *CancelAskOrder) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
