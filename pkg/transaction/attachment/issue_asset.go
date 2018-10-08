package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type IssueAssetAttachment struct {
	NumName        uint8  `struct:"uint8,sizeof=Name" json:"-"`
	Name           string `json:"name"`
	NumDescription uint16 `struct:"uint16,sizeof=Description" json:"-"`
	Description    string `json:"description"`
	Quantity       uint64 `json:"quantityQNT"`
	Decimals       uint8  `json:"decimals"`
	Comment        string `struct:"-" json:"comment,omitempty"`
}

func IssueAssetAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment IssueAssetAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Name) + 2 + len(attachment.Description) + 8 + 1, err
}

func (attachment *IssueAssetAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
