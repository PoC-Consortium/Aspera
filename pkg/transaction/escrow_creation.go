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

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachmet.Amount)
	e.WriteUint32(tx.Attachmet.Deadline)
	e.WriteUint8(tx.Attachmet.DeadlineAction)
	e.WriteUint8(uint8(tx.Attachment.RequiredSigners))
	for _, signer := range tx.Attachment.Signers {
		e.WriteUint64(signer)
	}
}

func (tx *RewardRecipientAssignment) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 4 + 1 + 1 + len(tx.Attachment.Signers)*8
}
