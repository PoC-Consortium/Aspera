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

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint8(uint8(tx.Attachment.Decision))

	return e.Bytes()
}

func (tx *EscrowSign) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 1
}
