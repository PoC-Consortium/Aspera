package attachment

type SendMoneyAttachment struct{}

func SendMoneyAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	return &SendMoneyAttachment{}, 0, nil
}

func (attachment *SendMoneyAttachment) ToBytes() ([]byte, error) {
	return []byte{}, nil
}
