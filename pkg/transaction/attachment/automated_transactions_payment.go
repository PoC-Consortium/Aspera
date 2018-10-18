package attachment

type AutomatedTransactionsPayment struct{}

func (attachment *AutomatedTransactionsPayment) FromBytes(bs []byte, version uint8) (int, error) {
	return 0, nil
}

func (attachment *AutomatedTransactionsPayment) ToBytes(version uint8) ([]byte, error) {
	return []byte{}, nil
}

func (attachment *AutomatedTransactionsPayment) GetFlag() uint32 {
	return StandardAttachmentFlag
}
