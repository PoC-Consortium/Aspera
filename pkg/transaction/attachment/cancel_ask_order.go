package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type CancelAskOrderAttachment struct {
	Order uint64
}

func CancelAskOrderAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment CancelAskOrderAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8, err
}

func (attachment *CancelAskOrderAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
