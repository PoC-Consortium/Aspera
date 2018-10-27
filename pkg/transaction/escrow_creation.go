package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type EscrowCreation struct {
	*pb.EscrowCreation
}

func (tx *EscrowCreation) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Amount)
	e.WriteUint32(tx.Attachment.Deadline)
	e.WriteUint8(uint8(tx.Attachment.DeadlineAction))
	e.WriteUint8(uint8(tx.Attachment.RequiredSigners))
	for _, signer := range tx.Attachment.Signers {
		e.WriteUint64(signer)
	}

	return e.Bytes()
}

func (tx *EscrowCreation) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 4 + 1 + 1 + len(tx.Attachment.Signers)*8
}
