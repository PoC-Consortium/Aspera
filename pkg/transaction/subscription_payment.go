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

	WriteHeader(e, tx.Header)

	for _, id := range tx.Attachment.Ids {
		e.WriteUint64(id)
	}

	return e.Bytes()
}

func (tx *SubscriptionPayment) SizeInBytes() int {
	return HeaderSize(tx.Header) + len(tx.Attachment.Ids)*8
}
