package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type EscrowResultAttachment struct {
	EscrowID uint64
	Decision uint8
}

func EscrowResultAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment EscrowResultAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 1, err
}

func (attachment *EscrowResultAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
