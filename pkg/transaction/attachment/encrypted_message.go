package attachment

import (
	"bytes"
	"encoding/binary"

	"github.com/ac0v/aspera/pkg/parsing"
	"gopkg.in/restruct.v1"
)

type EncryptedMessage struct {
	IsText bool `struct:"-" json:"isText"`

	// IsText is encoded as a signle bit
	IsTextAndLen int32  `json:"-"`
	Data         []byte `json:"data"`
	Nonce        []byte `json:"nonce"`
}

type EncryptedMessageAttachment struct {
	Message *EncryptedMessage `json:"encryptedMessage"`
	Version int8              `struct:"-" json:"version.EncryptedMessage,omitempty"`
}

func (attachment *EncryptedMessageAttachment) ToBytes(version uint8) ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version > 0 {
		return append([]byte{version}, bs...), nil
	}

	return bs, nil
}

func EncryptedMessageAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var message EncryptedMessageAttachment

	r := bytes.NewReader(bs)

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return nil, 0, err
	}

	message.Message.IsTextAndLen = isTextAndLen
	message.Message.IsText = isText

	message.Message.Data = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &message.Message.Data); err != nil {
		return nil, 0, err
	}

	message.Message.Nonce = make([]byte, 32)
	err = binary.Read(r, binary.LittleEndian, &message.Message.Nonce)

	return &message, 4 + int(len), err
}
