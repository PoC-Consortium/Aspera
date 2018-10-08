package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type BuyAliasAttachment struct {
	NumAlias uint8  `struct:"uint8,sizeof=Alias" json:"-"`
	Alias    string `json:"alias"`
	Version  int8   `struct:"-" json:"version.AliasAssignment,omitempty"`
}

func BuyAliasAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment BuyAliasAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Alias), err
}

func (attachment *BuyAliasAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
