package attachment

import (
	"bytes"
	"encoding/binary"

	"gopkg.in/restruct.v1"

	"github.com/ac0v/aspera/pkg/parsing"
)

type EncryptedToSelfMessageAttachment struct {
	IsText bool `struct:"-"`

	// IsText is encoded as a signle bit
	IsTextAndLen int32
	Data         []byte
	Nonce        []byte
}

func (attachment *EncryptedToSelfMessageAttachment) ToBytes(version uint8) ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version > 0 {
		return append([]byte{version}, bs...), nil
	}

	return bs, nil
}

func EncryptedToSelfMessageAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var message EncryptedToSelfMessageAttachment

	r := bytes.NewReader(bs)

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return nil, 0, err
	}

	message.IsTextAndLen = isTextAndLen
	message.IsText = isText

	message.Data = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &message.Data); err != nil {
		return nil, 0, err
	}

	message.Nonce = make([]byte, 32)
	err = binary.Read(r, binary.LittleEndian, &message.Nonce)

	return &message, 4 + int(len), err
}
