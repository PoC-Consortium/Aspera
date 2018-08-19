package attachment

type SendMoneyTransaction struct{}

func SendMoneyTransactionFromBytes(bs []byte) (Attachment, int, error) {
	return &SendMoneyTransaction{}, 0, nil
}
