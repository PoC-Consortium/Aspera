package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type AdvancedPaymentSubscriptionPaymentAttachment struct {
	SubscriptionID uint64 `json:"subscriptionId"`
}

func AdvancedPaymentSubscriptionPaymentAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment BuyAliasAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Alias), err
}

func (attachment *AdvancedPaymentSubscriptionPaymentAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
