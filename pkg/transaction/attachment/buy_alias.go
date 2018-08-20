package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type BuyAliasAttachment struct {
	NumAlias uint8 `struct:"uint8,sizeof=Alias"`
	Alias    []byte
}

func BuyAliasAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment BuyAliasAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Alias), err
}

func (attachment *BuyAliasAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
