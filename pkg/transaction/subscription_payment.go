package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type SubscriptionPayment struct {
	*pb.SubscriptionPayment
}

func (tx *SubscriptionPayment) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	for _, id := range tx.Attachment.SubscriptionIDs {
		e.WriteUint64(id)
	}
}

func (tx *SubscriptionPayment) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + len(tx.Attachment.SubscriptionIDs)*8
}
