package attachment

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"gopkg.in/restruct.v1"
	"strconv"
)

type Payment struct {
	Recip  uint64
	Amount uint64
}

type SendMoneyMultiAttachment struct {
	NumRecipsAndAmounts uint8     `struct:"uint8,sizeof=RecipsAndAmounts" json:"-"`
	RecipsAndAmounts    []Payment `json:"recipients"`
	Version             int8      `struct:"-" json:"version.MultiOutCreation,omitempty"`
}

func SendMoneyMultiAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SendMoneyMultiAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.RecipsAndAmounts)*(8+8), err
}

func (attachment *SendMoneyMultiAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (p *Payment) UnmarshalJSON(b []byte) error {
	var v []uint64
	if err := json.Unmarshal(bytes.Replace(b, []byte(`"`), []byte(""), -1), &v); err != nil {
		return err
	}
	p.Recip = v[0]
	p.Amount = v[1]

	return nil
}

func (p *Payment) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{strconv.FormatUint(p.Recip, 10), strconv.FormatUint(p.Amount, 10)})
}
