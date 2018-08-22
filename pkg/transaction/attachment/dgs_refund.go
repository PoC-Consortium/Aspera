package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsRefundAttachment struct {
	Purchase  uint64
	RefundNQT uint64
}

func DgsRefundAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment DgsRefundAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8, err
}

func (attachment *DgsRefundAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
