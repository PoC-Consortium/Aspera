package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type EscrowSign struct {
	Escrow   uint64
	Decision uint8
}

func (attachment *EscrowSign) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 1, err
}

func (attachment *EscrowSign) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *EscrowSign) GetFlag() uint32 {
	return StandardAttachmentFlag
}
