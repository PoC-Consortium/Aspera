package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SendMoneyEscrow struct {
	AmountNQT      uint64 `json:"amountNQT,omitempty,string"`
	EscrowDeadline uint32
	DeadlineAction uint8 `json:"xdeadlineAction,omitempty,string"` // ToDo: map enum :-()
	NumSignees     uint8 `struct:"uint8,sizeof=Signees"`
	TotalSignees   uint8
	Signees        []uint64 // TODO: implement json marshaler
}

func (attachment *SendMoneyEscrow) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 4 + 1 + 1 + 1 + len(attachment.Signees)*8, err
}

func (attachment *SendMoneyEscrow) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
