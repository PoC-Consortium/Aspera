package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPurchaseAttachment struct {
	Goods                     uint64
	Quantity                  uint32
	PriceNQT                  uint64
	DeliveryDeadlineTimestamp uint32
}

func DgsPurchaseAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment DgsPurchaseAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 4 + 8 + 4, err
}

func (attachment *DgsPurchaseAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
