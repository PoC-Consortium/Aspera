package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsListingAttachment struct {
	NumName        uint16 `struct:"uint16,sizeof=Name"`
	Name           []byte
	NumDescription uint16 `struct:"uint16,sizeof=Description"`
	Description    []byte
	NumTags        uint16 `struct:"uint16,sizeof=Tags"`
	Tags           []byte
	Quantity       uint32
	PriceNQT       uint64
}

func DgsListingAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment DgsListingAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 2 + len(attachment.Name) + 2 + len(attachment.Description) + 2 + len(attachment.Tags) + 4 + 8, err
}

func (attachment *DgsListingAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
