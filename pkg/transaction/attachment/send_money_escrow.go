package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyEscrowAttachment struct {
	AmountNQT      uint64
	EscrowDeadline uint32
	DeadlineAction uint8
	NumSignees     uint8 `struct:"uint8,sizeof=Signees"`
	TotalSignees   uint8
	Signees        []uint64
}

func SendMoneyEscrowAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SendMoneyEscrowAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 8 + 4 + 1 + 1 + 1 + len(attachment.Signees)*8, err
}

func (attachment *SendMoneyEscrowAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
