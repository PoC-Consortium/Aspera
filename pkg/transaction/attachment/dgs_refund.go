package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsRefundAttachment struct {
	Purchase  uint64
	RefundNQT uint64
}

func DgsRefundAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment DgsRefundAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8, err
}

func (attachment *DgsRefundAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
