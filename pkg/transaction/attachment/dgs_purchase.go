package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPurchase struct {
	Goods                     uint64 `json:"goods,omitempty,string"`
	Quantity                  uint32 `json:"quantity,omitempty"`
	PriceNQT                  uint64 `json:"priceNQT,omitempty"`
	DeliveryDeadlineTimestamp uint32 `json:"deliveryDeadlineTimestamp,omitempty"`
	Version                   int8   `struct:"-" json:"version.DigitalGoodsPurchase,omitempty"`
}

func (attachment *DgsPurchase) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 4 + 8 + 4, err
}

func (attachment *DgsPurchase) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
