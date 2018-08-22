package attachment

type SendMoneyAttachment struct{}

func SendMoneyAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	return &SendMoneyAttachment{}, 0, nil
}

func (attachment *SendMoneyAttachment) ToBytes(version uint8) ([]byte, error) {
	return []byte{}, nil
}
