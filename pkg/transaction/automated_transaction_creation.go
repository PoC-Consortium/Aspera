package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	AutomatedTransactionsCreationType    = 22
	AutomatedTransactionsCreationSubType = 0
)

type AutomatedTransactionsCreation struct {
	*pb.AutomatedTransactionsCreation
}

func (tx *AutomatedTransactionsCreation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes([]byte(tx.Attachment.Description))
	e.WriteBytes(tx.Attachment.Bytes)
}

func (tx *AutomatedTransactionsCreation) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Description) + 2 + len(tx.Attachment.Description) + len(tx.Attachment.Bytes)
}

func (tx *AutomatedTransactionsCreation) GetType() uint16 {
	return AutomatedTransactionsCreationSubType<<8 | AutomatedTransactionsCreationType
}
