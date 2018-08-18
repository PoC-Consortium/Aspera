package transaction

type SendMoneyTransaction struct{}

func SendMoneyTransactionFromBytes(bs []byte) (Transaction, error) {
	return &SendMoneyTransaction{}, nil
}
