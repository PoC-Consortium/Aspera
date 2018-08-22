package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SubscriptionCancelAttachment struct {
	Subscription uint64
}

func SubscriptionCancelAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SubscriptionCancelAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8, err
}

func (attachment *SubscriptionCancelAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
