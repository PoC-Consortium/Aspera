package attachment

type SetRewardRecipientAttachment struct{}

func SetRewardRecipientAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment SetRewardRecipientAttachment
	return &attachment, 0, nil
}

func (attachment *SetRewardRecipientAttachment) ToBytes(version uint8) ([]byte, error) {
	return []byte{}, nil
}
