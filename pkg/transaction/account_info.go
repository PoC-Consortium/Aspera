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

	WriteHeader(e, tx.Header)

	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes([]byte(tx.Attachment.Description))

	return e.Bytes()
}

func (tx *AccountInfo) SizeInBytes() int {
	return HeaderSize(tx.Header) + 1 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description)
}
