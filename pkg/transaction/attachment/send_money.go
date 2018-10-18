package attachment

type SendMoney struct{}

func (attachment *SendMoney) FromBytes(bs []byte, version uint8) (int, error) {
	return 0, nil
}

func (attachment *SendMoney) ToBytes(version uint8) ([]byte, error) {
	return []byte{}, nil
}

func (attachment *SendMoney) GetFlag() uint32 {
	return StandardAttachmentFlag
}
