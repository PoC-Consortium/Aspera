package attachment

import (
	"encoding/binary"
	"fmt"

	"gopkg.in/restruct.v1"
)

type EscrowSign struct {
	Escrow       uint64 `json:"escrow"`
	DecisionByte uint8  `json:"-"`
	Decision     string `struct:"-" json:"decision"`
}

func (attachment *EscrowSign) FromBytes(bs []byte, version uint8) (int, error) {
	if err := restruct.Unpack(bs, binary.LittleEndian, attachment); err != nil {
		return 0, err
	}
	if decision, exists := escrowDeadlineActionNameOf[attachment.DecisionByte]; exists {
		attachment.Decision = decision
	} else {
		return 0, fmt.Errorf("unknown escrow decision byte: %d", attachment.DecisionByte)
	}
	return 8 + 1, nil
}

func (attachment *EscrowSign) ToBytes(version uint8) ([]byte, error) {
	decisionByte, exists := escrowDeadlineActionIdOf[attachment.Decision]
	if !exists {
		return nil, fmt.Errorf("unknown escrow decision: %s", attachment.Decision)
	}
	attachment.DecisionByte = decisionByte
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *EscrowSign) GetFlag() uint32 {
	return StandardAttachmentFlag
}
