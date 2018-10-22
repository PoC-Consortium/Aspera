package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type EscrowSign struct {
	*pb.EscrowSign
}

func (tx *EscrowSign) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachmet.EscrowID)
	e.WriteUint8(tx.Attachmet.Decision)
}

func (tx *RewardRecipientAssignment) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 1
}
