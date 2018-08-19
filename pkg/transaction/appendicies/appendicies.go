package appendicies

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/ac0v/aspera/pkg/parsing"
)

const (
	maxInt32      = 2147483647
	maxMessageLen = 1000
)

var errMessageTooLong = errors.New("message too long")

type Message struct {
	IsText  bool
	Len     int32
	Content []byte
}

type EncryptedMessage struct {
	IsText bool
	Len    int32
	Data   []byte
	Nonce  []byte
}

type PublicKeyAnnouncement struct {
	PublicKey []byte
}

type EncryptedToSelfMessage struct {
	IsText bool
	Len    int32
	Data   []byte
	Nonce  []byte
}

type Appendices struct {
	Message                *Message
	EncryptedMessage       *EncryptedMessage
	PublicKeyAnnouncement  *PublicKeyAnnouncement
	EncryptedToSelfMessage *EncryptedMessage
}

func getMessageLengthAndType(r io.Reader) (int32, bool, error) {
	var len int32
	var isText bool

	if err := binary.Read(r, binary.LittleEndian, &len); err != nil {
		return len, isText, err
	}

	isText = len < 0
	if isText {
		len &= maxInt32
	}

	if len > maxMessageLen {
		return len, isText, errMessageTooLong
	}

	return len, isText, nil
}

func messageFromBytes(r io.Reader) (*Message, error) {
	var message Message

	len, isText, err := getMessageLengthAndType(r)
	if err != nil {
		return nil, err
	}

	message.Len = len
	message.IsText = isText

	message.Content = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &message.Content); err != nil {
		return nil, err
	}

	return &message, nil
}

func encryptedMessageFromBytes(r io.Reader) (*EncryptedMessage, error) {
	var message EncryptedMessage

	len, isText, err := getMessageLengthAndType(r)
	if err != nil {
		return nil, err
	}

	message.Len = len
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
		m, err := messageFromBytes(r)
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
