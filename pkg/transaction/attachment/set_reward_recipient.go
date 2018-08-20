package attachment

type SetRewardRecipientAttachment struct{}

func SetRewardRecipientAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment SetRewardRecipientAttachment
	return &attachment, 0, nil
}

func (attachment *SetRewardRecipientAttachment) ToBytes() ([]byte, error) {
	return []byte{}, nil
}
