package attachment

import (
	"bytes"
	"encoding/binary"

	"github.com/ac0v/aspera/pkg/transaction/appendicies"
	"gopkg.in/restruct.v1"
)

type SendMessageAttachment struct {
	*appendicies.Message
}

func SendMessageAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment SendMessageAttachment
	r := bytes.NewReader(bs)

	m, err := appendicies.MessageFromBytes(r)
	if err != nil {
		return nil, 0, err
	}

	attachment.Message = m

	return &attachment, int(r.Size()) - r.Len(), nil
}

func (attachment *SendMessageAttachment) ToBytes() ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment.Message)
}
