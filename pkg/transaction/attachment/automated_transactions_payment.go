package attachment

type AutomatedTransactionsPaymentAttachment struct{}

func AutomatedTransactionsPaymentAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	return new(AutomatedTransactionsPaymentAttachment), 0, nil
}

func (attachment *AutomatedTransactionsPaymentAttachment) ToBytes(version uint8) ([]byte, error) {
	return []byte{}, nil
}
