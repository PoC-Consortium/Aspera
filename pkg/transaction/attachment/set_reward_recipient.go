package attachment

type SetRewardRecipientTransaction struct{}

func SetRewardRecipientTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SetRewardRecipientTransaction
	return &tx, 0, nil
}
