package transaction

import ()

type AtPaymentTransaction struct{}

func AtPaymentTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx AtPaymentTransaction
	return &tx, 0, nil
}
