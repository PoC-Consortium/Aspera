package attachment

type DummyAttachment struct {
}

func (attachment *DummyAttachment) ToBytes(version uint8) ([]byte, error) {
	return nil, nil
}
