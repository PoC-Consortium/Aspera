package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsDelistingAttachment struct {
	Goods uint64
}

func DgsDelistingAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment DgsDelistingAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8, err
}

func (attachment *DgsDelistingAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
