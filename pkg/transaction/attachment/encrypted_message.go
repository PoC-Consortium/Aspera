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

func (attachment *EncryptedMessage) FromBytes(bs []byte, version uint8) (int, error) {
	r := bytes.NewReader(bs)

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return 0, err
	}

	attachment.IsTextAndLen = isTextAndLen
	attachment.IsText = isText

	attachment.Data = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &attachment.Data); err != nil {
		return 0, err
	}

	attachment.Nonce = make([]byte, 32)
	err = binary.Read(r, binary.LittleEndian, &attachment.Nonce)

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
