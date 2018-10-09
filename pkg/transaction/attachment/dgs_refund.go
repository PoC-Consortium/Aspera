package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsRefund struct {
	Purchase  uint64 `json:"purchase,omitempty,string"`
	RefundNQT uint64 `json:"refundNQT,omitempty"`
	Version   int8   `struct:"-" json:"version.DigitalGoodsRefund,omitempty"`
}

func (attachment *DgsRefund) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 8, err
}

func (attachment *DgsRefund) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
