package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type EffectiveBalanceLeasing struct {
	*pb.EffectiveBalanceLeasing
}

func (tx *EffectiveBalanceLeasing) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint32(tx.Attachment.Period)

	return e.Bytes()
}

func (tx *EffectiveBalanceLeasing) SizeInBytes() int {
	return HeaderSize(tx.Header) + 4
}
