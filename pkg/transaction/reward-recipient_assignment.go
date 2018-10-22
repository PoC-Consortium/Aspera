package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type RewardRecipientAssignment struct {
	*pb.RewardRecipientAssignment
}

func (tx *RewardRecipientAssignment) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)
}

func (tx *RewardRecipientAssignment) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader)
}
