package attachment

import (
	"bytes"
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PublicKeyAnnouncement struct {
	PublicKey []byte
}

func (attachment *PublicKeyAnnouncement) ToBytes(version uint8) ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version > 0 {
		return append([]byte{version}, bs...), nil
	}

	return bs, nil
}

func PublicKeyAnnouncementAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var message PublicKeyAnnouncement

	r := bytes.NewReader(bs)

	message.PublicKey = make([]byte, 32)
	err := binary.Read(r, binary.LittleEndian, &message.PublicKey)

	return &message, 32, err
}
