package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type CancelAskOrderAttachment struct {
	Order   uint64 `json:"order,omitempty,string"`
	Version int8   `struct:"-" json:"version.AskOrderCancellation,omitempty"`
}

func CancelAskOrderAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment CancelAskOrderAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8, err
}

func (attachment *CancelAskOrderAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
