package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AutomatedTransactionPayment struct {
	*pb.AutomatedTransactionPayment
}

func (tx *AutomatedTransactionPayment) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)
}

func (tx *AutomatedTransactionPayment) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader)
}
