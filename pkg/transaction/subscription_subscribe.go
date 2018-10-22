package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type SubscriptionSubscribe struct {
	*pb.SubscriptionSubscribe
}

func (tx *SubscriptionSubscribe) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint32(tx.Attachmet.Frequency)
}
func (tx *SubscriptionSubscribe) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 4
}
