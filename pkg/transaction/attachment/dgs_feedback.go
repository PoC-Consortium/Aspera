package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsFeedbackAttachment struct {
	Purchase uint64
}

func DgsFeedbackAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment DgsFeedbackAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8, err
}

func (attachment *DgsFeedbackAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
