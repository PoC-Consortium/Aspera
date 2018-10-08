package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsRefundAttachment struct {
	Purchase  uint64 `json:"purchase,omitempty,string"`
	RefundNQT uint64 `json:"refundNQT,omitempty"`
	Version   int8   `struct:"-" json:"version.DigitalGoodsRefund,omitempty"`
}

func DgsRefundAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment DgsRefundAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8, err
}

func (attachment *DgsRefundAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
