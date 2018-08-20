package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsDelistingAttachment struct {
	Goods uint64
}

func DgsDelistingAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment DgsDelistingAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8, err
}

func (attachment *DgsDelistingAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
