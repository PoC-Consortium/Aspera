package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneySubscriptionAttachment struct {
	Frequency uint32 `json:"frequency"`
	Version   int8   `struct:"-" json:"version.SubscriptionSubscribe,omitempty"`
}

func SendMoneySubscriptionAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SendMoneySubscriptionAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 4, err
}

func (attachment *SendMoneySubscriptionAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
