package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneySubscriptionAttachment struct {
	Frequency uint32
}

func SendMoneySubscriptionAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment SendMoneySubscriptionAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 4, err
}

func (attachment *SendMoneySubscriptionAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
