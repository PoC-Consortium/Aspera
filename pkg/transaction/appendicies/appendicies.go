package appendicies

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/ac0v/aspera/pkg/parsing"
	"gopkg.in/restruct.v1"
)

const (
	maxInt32      = 2147483647
	maxMessageLen = 1000
)

var errMessageTooLong = errors.New("message too long")

type Message struct {
	IsText bool `struct:"-" json:"messageIsText"`

	// IsText is encoded as a single bit
	IsTextAndLen int32  `json:"-"`
	Content      string `json:"message"`
	Version      int8   `struct:"-" json:"version.Message,omitempty"`
}

type EncryptedMessage struct {
	IsText bool `struct:"-" json:"messageIsText"`

	// IsText is encoded as a signle bit
	IsTextAndLen int32
	Data         []byte
	Nonce        []byte
}

type PublicKeyAnnouncement struct {
	PublicKey []byte
}

type EncryptedToSelfMessage struct {
	IsText bool `struct:"-"`

	// IsText is encoded as a signle bit
	IsTextAndLen int32
	Data         []byte
	Nonce        []byte
}

type Appendices struct {
	Message                *Message
	EncryptedMessage       *EncryptedMessage
	PublicKeyAnnouncement  *PublicKeyAnnouncement
	EncryptedToSelfMessage *EncryptedMessage
}

func getMessageLengthAndType(r io.Reader) (int32, int32, bool, error) {
	var len, isTextAndLen int32
	var isText bool

	if err := binary.Read(r, binary.LittleEndian, &isTextAndLen); err != nil {
		return 0, 0, isText, err
	}

	isText = isTextAndLen < 0
	if isText {
		len = isTextAndLen & maxInt32
	} else {
		len = isTextAndLen
	}

	if len > maxMessageLen {
		return 0, 0, isText, errMessageTooLong
	}

	return len, isTextAndLen, isText, nil
}

func MessageFromBytes(r io.Reader) (*Message, error) {
	var message Message

	len, isTextAndLen, isText, err := getMessageLengthAndType(r)
	if err != nil {
		return nil, err
	}

	message.IsTextAndLen = isTextAndLen
	message.IsText = isText

	content := make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &content); err != nil {
		return nil, err
	}
	message.Content = string(content)

	return &message, nil
}

func encryptedMessageFromBytes(r io.Reader) (*EncryptedMessage, error) {
	var message EncryptedMessage

	len, isTextAndLen, isText, err := getMessageLengthAndType(r)
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

func publicKeyAnnouncementFromBytes(r io.Reader) (*PublicKeyAnnouncement, error) {
	var message PublicKeyAnnouncement

	message.PublicKey = make([]byte, 32)
	err := binary.Read(r, binary.LittleEndian, &message.PublicKey)

	return &message, err
}

func FromBytes(bs []byte, flags uint32, version uint8) (*Appendices, error) {
	var appendicies Appendices

	r := bytes.NewReader(bs)

	if flags&(1<<0) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := MessageFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.Message = m
	}

	if flags&(1<<1) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := encryptedMessageFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.EncryptedMessage = m
	}

	if flags&(1<<2) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := publicKeyAnnouncementFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.PublicKeyAnnouncement = m
	}

	if flags&(1<<3) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := encryptedMessageFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.EncryptedToSelfMessage = m
	}

	return &appendicies, nil
}

func (appendicies *Appendices) ToBytes(version uint8) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if appendicies.Message != nil {
		if version > 0 {
			if err := binary.Write(buf, binary.LittleEndian, version); err != nil {
				return nil, err
			}
		}

		bs, err := restruct.Pack(binary.LittleEndian, appendicies.Message)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(bs); err != nil {
			return nil, err
		}
	}

	if appendicies.EncryptedMessage != nil {
		if version > 0 {
			if err := binary.Write(buf, binary.LittleEndian, version); err != nil {
				return nil, err
			}
		}

		bs, err := restruct.Pack(binary.LittleEndian, appendicies.EncryptedMessage)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(bs); err != nil {
			return nil, err
		}
	}

	if appendicies.PublicKeyAnnouncement != nil {
		if version > 0 {
			if err := binary.Write(buf, binary.LittleEndian, version); err != nil {
				return nil, err
			}
		}

		bs, err := restruct.Pack(binary.LittleEndian, appendicies.PublicKeyAnnouncement)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(bs); err != nil {
			return nil, err
		}
	}

	if appendicies.EncryptedToSelfMessage != nil {
		if version > 0 {
			if err := binary.Write(buf, binary.LittleEndian, version); err != nil {
				return nil, err
			}
		}

		bs, err := restruct.Pack(binary.LittleEndian, appendicies.EncryptedToSelfMessage)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(bs); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
