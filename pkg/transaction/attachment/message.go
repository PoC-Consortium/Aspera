package attachment

import (
	"bytes"
	"encoding/binary"

	"github.com/ac0v/aspera/pkg/parsing"
	"gopkg.in/restruct.v1"
)

type MessageAttachment struct {
	IsText bool `struct:"-" json:"messageIsText"`

	// IsText is encoded as a single bit
	IsTextAndLen int32  `json:"-"`
	Content      string `json:"message"`
	Version      int8   `struct:"-" json:"version.Message,omitempty"`
}

func (attachment *MessageAttachment) ToBytes(version uint8) ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version > 0 {
		return append([]byte{version}, bs...), nil
	}

	return bs, nil
}

func MessageAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var message MessageAttachment

	r := bytes.NewReader(bs)

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return nil, 0, err
	}

	message.IsTextAndLen = isTextAndLen
	message.IsText = isText

	content := make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &content); err != nil {
		return nil, 0, err
	}
	message.Content = string(content)

	return &message, 4 + int(len), nil
}
