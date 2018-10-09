package attachment

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"gopkg.in/restruct.v1"
)

type Recipient struct {
	Recip uint64
}

type SendMoneyMultiSameAttachment struct {
	RecipCount uint8       `struct:"uint8,sizeof=Recips" json:"-"`
	Recips     []Recipient `json:"recipients"`
	Version    int8        `struct:"-" json:"version.MultiSameOutCreation,omitempty"`
}

func SendMoneyMultiSameAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SendMoneyMultiSameAttachment
	err := restruct.Unpack(bs, binary.LittleEndian, &attachment)
	return &attachment, 1 + len(attachment.Recips)*8, err
}

func (attachment *SendMoneyMultiSameAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (r *Recipient) UnmarshalJSON(b []byte) error {
	var err error
	r.Recip, err = strconv.ParseUint(string(bytes.Replace(b, []byte(`"`), []byte(""), 2)), 10, 64)
	return err
}

func (p *Recipient) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(p.Recip, 10) + `"`), nil
}
