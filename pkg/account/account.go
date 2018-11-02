package account

import (
	"errors"

	"github.com/ac0v/aspera/pkg/account/pb"
	"github.com/ac0v/aspera/pkg/crypto"
	"github.com/ac0v/aspera/pkg/crypto/rsencoding"

	"github.com/golang/protobuf/proto"
)

var (
	ErrPublicKeyInvalidLen = errors.New("public key has invalid length")
)

type Account struct {
	*pb.Account
}

func NewAccount(publicKey []byte, balance int64) *Account {
	id := publicKeyToID(publicKey)
	return &Account{
		Account: &pb.Account{
			Id:              id,
			PublicKey:       publicKey,
			Balance:         balance,
			RewardRecipient: id,
			Address:         rsencoding.Encode(id),
		},
	}
}

func (a *Account) ToBytes() []byte {
	if bs, err := proto.Marshal(a.Account); err == nil {
		return bs
	} else {
		panic(err)
	}
}

func FromBytes(bs []byte) *Account {
	var a pb.Account
	if err := proto.Unmarshal(bs, &a); err == nil {
		return &Account{Account: &a}
	} else {
		panic(err)
	}
}

func publicKeyToID(publicKey []byte) uint64 {
	_, id := crypto.BytesToHashAndID(publicKey)
	return id
}
