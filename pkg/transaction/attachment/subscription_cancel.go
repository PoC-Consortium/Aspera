package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SubscriptionCancel struct {
	Subscription uint64 `json:"subscriptionId,string"`
	Version      int8   `struct:"-" json:"version.SubscriptionCancel,omitempty"`
}

func (attachment *SubscriptionCancel) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8, err
}

func (attachment *SubscriptionCancel) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
