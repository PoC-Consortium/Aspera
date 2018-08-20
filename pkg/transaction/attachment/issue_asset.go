package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type IssueAssetAttachment struct {
	NumName        uint8 `struct:"uint8,sizeof=Name"`
	Name           []byte
	NumDescription uint16 `struct:"uint16,sizeof=Description"`
	Description    []byte
	Quantity       uint64
	Decimals       uint8
}

func IssueAssetAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment IssueAssetAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Name) + 2 + len(attachment.Description) + 8 + 1, err
}

func (attachment *IssueAssetAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
