package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyMultiSameAttachment struct {
	RecipCount uint8 `struct:"uint8,sizeof=Recips"`
	Recips     []uint64
}

func SendMoneyMultiSameAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment SendMoneyMultiSameAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Recips)*8, err
}

func (attachment *SendMoneyMultiSameAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
