package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SetAliasAttachment struct {
	NumAliasName uint8 `struct:"uint8,sizeof=AliasName"`
	AliasName    []byte
	NumAliasURI  uint16 `struct:"uint16,sizeof=AliasURI"`
	AliasURI     []byte
}

func SetAliasAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment SetAliasAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.AliasName) + 2 + len(attachment.AliasURI), err
}

func (attachment *SetAliasAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
