package transaction

import ()

type AtPaymentTransaction struct{}

func AtPaymentTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx AtPaymentTransaction
	return &tx, nil
}
