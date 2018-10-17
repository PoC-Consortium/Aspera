package attachment

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/ac0v/aspera/pkg/parsing"
	"github.com/json-iterator/go"
)

var js = jsoniter.ConfigCompatibleWithStandardLibrary

type Message struct {
	IsText *bool `struct:"-" json:"messageIsText,omitempty"`

	// IsText is encoded as a single bit
	IsTextAndLen int32  `json:"-"`
	Content      string `json:"message,omitempty"`
	Version      int8   `struct:"-" json:"version.Message,omitempty"`
}

func (m *Message) UnmarshalJSON(bs []byte) error {
	type Alias Message
	ma := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := js.Unmarshal(bs, ma); err != nil {
		return err
	}

	if m.IsText != nil && *m.IsText {
		m.IsTextAndLen = int32(-len(m.Content))
	} else {
		m.IsTextAndLen = int32(len(m.Content) / 2)
	}

	return nil
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
	if isText {
		attachment.Content = string(content)
	} else {
		attachment.Content = hex.EncodeToString(content)
	}

	return 4 + int(len), nil
}

func (attachment *Message) ToBytes(version uint8) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if err := binary.Write(buf, binary.LittleEndian, attachment.IsTextAndLen); err != nil {
		return nil, err
	}

	if *attachment.IsText {
		if err := binary.Write(buf, binary.LittleEndian, []byte(attachment.Content)); err != nil {
			return nil, err
		}
	} else {
		content, err := hex.DecodeString(attachment.Content)
		if err != nil {
			return nil, err
		}

		if err = binary.Write(buf, binary.LittleEndian, content); err != nil {
			return nil, err
		}
	}

	if version > 0 {
		return append([]byte{version}, buf.Bytes()...), nil
	}

	return buf.Bytes(), nil
}
