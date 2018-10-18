package attachment

type SetRewardRecipient struct {
	Version int8 `struct:"-" json:"version.RewardRecipientAssignment,omitempty"`
}

func (attachment *SetRewardRecipient) FromBytes(bs []byte, version uint8) (int, error) {
	return 0, nil
}

func (attachment *SetRewardRecipient) ToBytes(version uint8) ([]byte, error) {
	return []byte{}, nil
}

func (attachment *SetRewardRecipient) GetFlag() uint32 {
	return StandardAttachmentFlag
}
