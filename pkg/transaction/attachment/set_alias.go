package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SetAliasAttachment struct {
	NumAliasName uint8  `struct:"uint8,sizeof=AliasName" json:"-"`
	AliasName    string `json:"alias"`
	NumAliasURI  uint16 `struct:"uint16,sizeof=AliasURI" json:"-"`
	AliasURI     string `json:"uri"`
	Version      int8   `struct:"-" json:"version.AliasAssignment,omitempty"`
}

func SetAliasAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SetAliasAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.AliasName) + 2 + len(attachment.AliasURI), err
}

func (attachment *SetAliasAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
