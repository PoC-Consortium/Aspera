package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsPriceChangeAttachment struct {
	Goods    uint64
	PriceNQT uint64
}

func DgsPriceChangeAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment DgsPriceChangeAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 8, err
}

func (attachment *DgsPriceChangeAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
