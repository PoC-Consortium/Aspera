package transaction

type sendMessageTransaction struct{}

func SendMessageTransactionFromBytes(bs []byte) (Transaction, error) {
	return &sendMessageTransaction{}, nil
}
