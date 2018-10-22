package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AccountInfo struct {
	*pb.AccountInfo
}

func (tx *AccountInfo) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint8(uint8(len(tx.Attachment.Alias)))
	e.WriteBytes([]byte(tx.Attachment.Alias))
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *AccountInfo) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 1 + len(tx.Attachment.Alias) + 8
}
