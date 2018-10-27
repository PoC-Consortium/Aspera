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

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Id)

	return e.Bytes()
}

func (tx *SubscriptionCancel) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8
}
