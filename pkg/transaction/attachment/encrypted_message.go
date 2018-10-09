package attachment

import (
	"bytes"
	"encoding/binary"

	"github.com/ac0v/aspera/pkg/parsing"
	"gopkg.in/restruct.v1"
)

type EncryptedMessage struct {
	IsText bool `struct:"-" json:"messageIsText"`

	// IsText is encoded as a signle bit
	IsTextAndLen int32
	Data         []byte
	Nonce        []byte
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

func EncryptedMessageFromBytes(r *bytes.Reader, version uint8) (*EncryptedMessage, error) {
	var message EncryptedMessage

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return nil, err
	}

	message.IsTextAndLen = isTextAndLen
	message.IsText = isText

	message.Data = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &message.Data); err != nil {
		return nil, err
	}

	message.Nonce = make([]byte, 32)
	err = binary.Read(r, binary.LittleEndian, &message.Nonce)

	return &message, err
}
