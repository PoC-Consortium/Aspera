package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AutomatedTransactionCreation struct {
	*pb.SubscriptionPayment
}

func (tx *AutomatedTransactionCreation) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteUint16(uint8(len(tx.Attachment.Description)))
	e.WriteBytes([]byte(tx.Attachment.Description))
	e.WriteBytes(tx.Attachment.CreationBytes)
}

func (tx *AutomatedTransactionCreation) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 1 + len(tx.Attachment.Description) +
		2 + len(tx.Attachment.Description) + len(tx.Attachment.CreationBytes)
}
