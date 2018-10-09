package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneySubscription struct {
	Frequency uint32 `json:"frequency"`
	Version   int8   `struct:"-" json:"version.SubscriptionSubscribe,omitempty"`
}

func (attachment *SendMoneySubscription) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 4, err
}

func (attachment *SendMoneySubscription) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
