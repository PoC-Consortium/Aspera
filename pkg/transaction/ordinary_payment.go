package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type OrdinaryPayment struct {
	*pb.OrdinaryPayment
}

func (tx *OrdinaryPayment) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	return e.Bytes()
}

func (tx *OrdinaryPayment) SizeInBytes() int {
	return HeaderSize(tx.Header)
}
