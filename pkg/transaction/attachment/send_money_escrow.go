package attachment

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"

	"gopkg.in/restruct.v1"
)

type Signer struct {
	Id uint64
}

type EscrowDeadlineAction struct {
	DeadlineAction uint8
}

type SendMoneyEscrow struct {
	AmountNQT       uint64               `json:"amountNQT,omitempty,string"`
	Deadline        uint32               `json:"deadline"`
	EscrowDeadline  EscrowDeadlineAction `json:"deadlineAction,omitempty"`
	RequiredSigners int8                 `json:"requiredSigners"`
	NumSignees      uint8                `struct:"uint8,sizeof=Signees" json:"-"`
	Signees         []Signer             `json:"signers"`
	Version         int8                 `struct:"-" json:"version.EscrowCreation,omitempty"`
}

func (attachment *SendMoneyEscrow) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 8 + 4 + 1 + 1 + 1 + len(attachment.Signees)*8, err
}

func (attachment *SendMoneyEscrow) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (s *Signer) UnmarshalJSON(b []byte) error {
	var err error
	s.Id, err = strconv.ParseUint(string(bytes.Replace(b, []byte(`"`), []byte(""), 2)), 10, 64)
	return err
}

func (s *Signer) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(s.Id, 10) + `"`), nil
}

func (e *EscrowDeadlineAction) UnmarshalJSON(b []byte) error {
	var exists bool
	if e.DeadlineAction, exists = escrowDeadlineActionIdOf[string(string(bytes.Replace(b, []byte(`"`), []byte(""), 2)))]; exists {
		return nil
	}

	return errors.New("unknown deadlineAction")
}

func (e *EscrowDeadlineAction) MarshalJSON() ([]byte, error) {
	if v, exists := escrowDeadlineActionNameOf[e.DeadlineAction]; exists {
		return []byte(`"` + v + `"`), nil
	}
	return nil, errors.New("unknown deadlineAction")
}

func (attachment *SendMoneyEscrow) GetFlag() uint32 {
	return StandardAttachmentFlag
}
