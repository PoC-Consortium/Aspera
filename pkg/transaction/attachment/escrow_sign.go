package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type EscrowSignAttachment struct {
	Escrow   uint64
	Decision uint8
}

func EscrowSignAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment EscrowSignAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 1, err
}

func (attachment *EscrowSignAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
