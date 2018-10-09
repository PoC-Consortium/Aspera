package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type AdvancedPaymentSubscriptionPayment struct {
	SubscriptionID uint64 `json:"subscriptionId"`
}

func (attachment *AdvancedPaymentSubscriptionPayment) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8, err
}

func (attachment *AdvancedPaymentSubscriptionPayment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
