package transaction

type SendMessageTransaction struct{}

func SendMessageTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SendMessageTransaction
	return &tx, 0, nil
}
