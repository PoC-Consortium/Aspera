package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsFeedback struct {
	Purchase uint64 `json:"purchase,omitempty,string"`
}

func (attachment *DgsFeedback) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8, err
}

func (attachment *DgsFeedback) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
