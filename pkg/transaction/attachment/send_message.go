package attachment

type SendMessageAttachment struct{}

func SendMessageAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment SendMessageAttachment
	return &attachment, 0, nil
}

func (attachment *SendMessageAttachment) ToBytes() ([]byte, error) {
	return []byte{}, nil
}
