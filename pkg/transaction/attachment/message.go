package attachment

import (
	"bytes"
	"encoding/binary"

	"github.com/ac0v/aspera/pkg/parsing"
	"gopkg.in/restruct.v1"
)

type Message struct {
	IsText *bool `struct:"-" json:"messageIsText,omitempty"`

	// IsText is encoded as a single bit
	IsTextAndLen int32  `json:"-"`
	Content      string `json:"message,omitempty"`
	Version      int8   `struct:"-" json:"version.Message,omitempty"`
}

func (attachment *Message) FromBytes(bs []byte, version uint8) (int, error) {
	r := bytes.NewReader(bs)

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return 0, err
	}

	attachment.IsTextAndLen = isTextAndLen
	attachment.IsText = &isText

	content := make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &content); err != nil {
		return 0, err
	}
	attachment.Content = string(content)

	return 4 + int(len), nil
}

func (attachment *Message) ToBytes(version uint8) ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version > 0 {
		return append([]byte{version}, bs...), nil
	}

	return bs, nil
}
