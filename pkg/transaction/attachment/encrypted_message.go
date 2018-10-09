package attachment

import (
	"bytes"
	"encoding/binary"

	"github.com/ac0v/aspera/pkg/parsing"
	"gopkg.in/restruct.v1"
)

type EncryptedMessageContainer struct {
	IsText bool `struct:"-" json:"isText"`

	// IsText is encoded as a single bit
	IsTextAndLen int32  `json:"-"`
	Data         []byte `json:"data"`
	Nonce        []byte `json:"nonce"`
}

type EncryptedMessage struct {
	Message EncryptedMessageContainer `json:"encryptedMessage"`
	Version int8                      `struct:"-" json:"version.EncryptedMessage,omitempty"`
}

func (attachment *EncryptedMessage) FromBytes(bs []byte, version uint8) (int, error) {
	r := bytes.NewReader(bs)

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return 0, err
	}

	attachment.Message.IsTextAndLen = isTextAndLen
	attachment.Message.IsText = isText

	attachment.Message.Data = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &attachment.Message.Data); err != nil {
		return 0, err
	}

	attachment.Message.Nonce = make([]byte, 32)
	err = binary.Read(r, binary.LittleEndian, &attachment.Message.Nonce)

	return 4 + int(len), err
}

func (attachment *EncryptedMessage) ToBytes(version uint8) ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version > 0 {
		return append([]byte{version}, bs...), nil
	}

	return bs, nil
}
