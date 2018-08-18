package transaction

type SetRewardRecipientTransaction struct{}

func SetRewardRecipientTransactionFromBytes(bs []byte) (Transaction, error) {
	var tx SetRewardRecipientTransaction
	return &tx, nil
}
