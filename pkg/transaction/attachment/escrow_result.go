package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type EscrowResult struct {
	EscrowID uint64 `json:"escrowId,string"`
	Decision uint8  `json:"decision"`
}

func (attachment *EscrowResult) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 1, err
}

func (attachment *EscrowResult) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *EscrowResult) GetFlag() uint32 {
	return StandardAttachmentFlag
}
