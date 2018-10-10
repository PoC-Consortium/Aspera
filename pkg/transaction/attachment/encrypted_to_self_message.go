package attachment

import (
	"bytes"
	"encoding/binary"

	"gopkg.in/restruct.v1"

	"github.com/ac0v/aspera/pkg/parsing"
)

type EncryptedToSelfMessageContainer struct {
	IsText bool `struct:"-" json:"isText"`

	// IsText is encoded as a signle bit
	IsTextAndLen int32  `json:"-"`
	Data         []byte `json:"data"`
	Nonce        []byte `json:"nonce"`
}

type EncryptedToSelfMessage struct {
	Message *EncryptedToSelfMessageContainer `json:"encryptToSelfMessage"`
	Version int8                    `struct:"-" json:"version.EncryptToSelfMessage,omitempty"`
}

func (attachment *EncryptedToSelfMessage) FromBytes(bs []byte, version uint8) (int, error) {
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

func (attachment *EncryptedToSelfMessage) ToBytes(version uint8) ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version > 0 {
		return append([]byte{version}, bs...), nil
	}

	return bs, nil
}
