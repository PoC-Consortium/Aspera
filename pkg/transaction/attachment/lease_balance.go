package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type LeaseBalanceAttachment struct {
	Period uint16
}

func LeaseBalanceAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment LeaseBalanceAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 2, err
}

func (attachment *LeaseBalanceAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
