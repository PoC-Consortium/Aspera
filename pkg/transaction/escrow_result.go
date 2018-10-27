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

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint8(uint8(tx.Attachment.Decision))

	return e.Bytes()
}

func (tx *EscrowResult) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 1
}
