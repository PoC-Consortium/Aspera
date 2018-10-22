package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type SubscriptionCancel struct {
	*pb.SubscriptionCancel
}

func (tx *SubscriptionCancel) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachmet.SubscriptionID)
}

func (tx *SubscriptionSubscribe) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8
}
