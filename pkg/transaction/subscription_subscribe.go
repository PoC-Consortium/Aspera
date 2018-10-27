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

	WriteHeader(e, tx.Header)

	e.WriteUint32(tx.Attachment.Frequency)

	return e.Bytes()
}
func (tx *SubscriptionSubscribe) SizeInBytes() int {
	return HeaderSize(tx.Header) + 4
}
