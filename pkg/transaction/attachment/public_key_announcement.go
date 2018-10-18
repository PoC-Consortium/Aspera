package attachment

import (
	"bytes"
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type PublicKeyAnnouncement struct {
	PublicKey []byte
}

func (attachment *PublicKeyAnnouncement) FromBytes(bs []byte, version uint8) (int, error) {
	r := bytes.NewReader(bs)

	attachment.PublicKey = make([]byte, 32)
	err := binary.Read(r, binary.LittleEndian, &attachment.PublicKey)

	return 32, err
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

func (attachment *PublicKeyAnnouncement) GetFlag() uint32 {
	return PublicKeyAnnouncementFlag
}
