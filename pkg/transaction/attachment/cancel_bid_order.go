package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type CancelBidOrder struct {
	Order   uint64 `json:"order,omitempty,string"`
	Version int8   `struct:"-" json:"version.BidOrderCancellation,omitempty"`
}

func (attachment *CancelBidOrder) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8, err
}

func (attachment *CancelBidOrder) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *CancelBidOrder) GetFlag() uint32 {
	return StandardAttachmentFlag
}
