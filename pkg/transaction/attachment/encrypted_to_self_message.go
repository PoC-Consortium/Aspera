package attachment

import (
	"bytes"
	"encoding/binary"
	"math"

	"gopkg.in/restruct.v1"

	jutils "github.com/ac0v/aspera/pkg/json"
	"github.com/ac0v/aspera/pkg/parsing"
)

type EncryptedToSelfMessageContainer struct {
	IsText bool `struct:"-" json:"isText"`

	// IsText is encoded as a signle bit
	IsTextAndLen int32           `json:"-"`
	Data         jutils.HexSlice `json:"data"`
	Nonce        jutils.HexSlice `json:"nonce"`
}

type EncryptedToSelfMessage struct {
	Message EncryptedToSelfMessageContainer `json:"encryptToSelfMessage"`
	Version int8                            `struct:"-" json:"version.EncryptToSelfMessage,omitempty"`
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
	attachment.Message.IsTextAndLen = int32(len(attachment.Message.Data))
	if attachment.Message.IsText {
		attachment.Message.IsTextAndLen |= math.MinInt32
	} else {
		// hex encoding
		// attachment.Message.IsTextAndLen /= 2
	}

	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func (attachment *EncryptedToSelfMessage) GetFlag() uint32 {
	return EncryptedToSelfMessageFlag
}
