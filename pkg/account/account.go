package account

import (
	"github.com/ac0v/aspera/pkg/encoding"
)

type Account struct {
	Id        uint64
	Balance   int64
	PublicKey []byte
}

func NewAccount(id uint64, balance int64) *Account {
	// TODO: calc public key
	return &Account{
		Id:      id,
		Balance: balance,
	}
}

func (a *Account) ToBytes() []byte {
	e := encoding.NewEncoder(a.SizeInBytes())
	e.WriteUint64(a.Id)
	e.WriteInt64(a.Balance)
	e.WriteBytes(a.PublicKey)
	return e.Bytes()
}

func FromBytes(bs []byte) *Account {
	d := encoding.NewDecoder(bs)
	id := d.ReadUint64()
	balance := d.ReadInt64()
	publicKey := d.ReadBytes(32)
	return &Account{
		Id:        id,
		Balance:   balance,
		PublicKey: publicKey,
	}
}

func (a *Account) SizeInBytes() int {
	return 8 + 8 + 32
}
