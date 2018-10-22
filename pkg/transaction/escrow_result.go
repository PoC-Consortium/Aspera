package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type EscrowResult struct {
	*pb.EscrowResult
}

func (tx *EscrowResult) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachmet.EscrowID)
	e.WriteUint8(tx.Attachmet.Decision)
}

func (tx *EscrowResult) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 1
}
