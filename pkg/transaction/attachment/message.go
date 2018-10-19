package attachment

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"

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
	attachment.IsTextAndLen = int32(len(attachment.Content))
	if attachment.IsText != nil && *attachment.IsText {
		attachment.IsTextAndLen |= math.MinInt32
	} else {
		// hex encoding
		attachment.IsTextAndLen /= 2
	}

	buf := bytes.NewBuffer(nil)

	if err := binary.Write(buf, binary.LittleEndian, attachment.IsTextAndLen); err != nil {
		return nil, err
	}

	if attachment.IsText != nil && *attachment.IsText {
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

	return buf.Bytes(), nil
}

func (attachment *Message) GetFlag() uint32 {
	return MessageFlag
}
