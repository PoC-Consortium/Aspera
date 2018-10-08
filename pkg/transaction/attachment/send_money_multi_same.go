package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyMultiSameAttachment struct {
	RecipCount uint8             `struct:"uint8,sizeof=Recips" json:"-"`
	Recips     UInt64StringSlice `json:"recipients"` // fix restruct
	Version    int8              `struct:"-" json:"version.MultiSameOutCreation,omitempty"`
}

func SendMoneyMultiSameAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SendMoneyMultiSameAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Recips)*8, err
}

func (attachment *SendMoneyMultiSameAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
