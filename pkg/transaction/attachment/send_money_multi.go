package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type Payment struct {
	Recip  uint64
	Amount uint64
}

type SendMoneyMultiAttachment struct {
	NumRecipsAndAmounts uint8 `struct:"uint8,sizeof=RecipsAndAmounts"`
	RecipsAndAmounts    []Payment
}

func SendMoneyMultiAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SendMoneyMultiAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.RecipsAndAmounts)*(8+8), err
}

func (attachment *SendMoneyMultiAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
