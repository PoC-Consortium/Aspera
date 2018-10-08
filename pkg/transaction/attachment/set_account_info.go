package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SetAccountInfoAttachment struct {
	NumName        uint8  `struct:"uint8,sizeof=Name" json:"-"`
	Name           string `json:"name"`
	NumDescription uint16 `struct:"uint16,sizeof=Description" json:"-"`
	Description    string `json:"description"`
	Version        int8   `struct:"-" json:"version.AccountInfo,omitempty"`
}

func SetAccountInfoAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SetAccountInfoAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Name) + 1 + len(attachment.Description), err
}

func (attachment *SetAccountInfoAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
