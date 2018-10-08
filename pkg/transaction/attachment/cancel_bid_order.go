package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type CancelBidOrderAttachment struct {
	Order   uint64 `json:"order,omitempty,string"`
	Version int8   `struct:"-" json:"version.BidOrderCancellation,omitempty"`
}

func CancelBidOrderAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment CancelBidOrderAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8, err
}

func (attachment *CancelBidOrderAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
