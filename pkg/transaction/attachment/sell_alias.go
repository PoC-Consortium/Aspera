package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SellAliasAttachment struct {
	NumAlias uint8 `struct:"uint8,sizeof=Alias"`
	Alias    []byte
	PriceNQT int64
}

func SellAliasAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SellAliasAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Alias) + 8, err
}

func (attachment *SellAliasAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
