package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AskOrderCancellation struct {
	*pb.AskOrderCancellation
}

func (tx *AskOrderCancellation) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Order)

	return e.Bytes()
}

func (tx *AskOrderCancellation) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8
}
