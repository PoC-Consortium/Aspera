package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type ArbitaryMessage struct {
	*pb.ArbitaryMessage
}

func (tx *ArbitaryMessage) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	return e.Bytes()
}

func (tx *ArbitaryMessage) SizeInBytes() int {
	return HeaderSize(tx.Header)
}
