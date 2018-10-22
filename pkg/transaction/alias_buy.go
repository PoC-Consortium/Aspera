package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AliasBuy struct {
	*pb.AliasBuy
}

func (tx *AliasBuy) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint8(uint8(len(tx.Attachment.Alias)))
	e.WriteBytes([]byte(tx.Attachment.Alias))
}

func (tx *AliasBuy) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 1 + len(tx.Attachment.Alias)
}
