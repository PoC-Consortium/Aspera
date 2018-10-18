package attachment

type Dummy struct {
}

func (attachment *Dummy) FromBytes(bs []byte, version uint8) (int, error) {
	return 0, nil
}

func (attachment *Dummy) ToBytes(version uint8) ([]byte, error) {
	return nil, nil
}

func (attachment *Dummy) GetFlag() uint32 {
	return StandardAttachmentFlag
}
