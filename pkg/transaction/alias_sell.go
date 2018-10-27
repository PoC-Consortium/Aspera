package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AliasSell struct {
	*pb.AliasSell
}

func (tx *AliasSell) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteInt64(tx.Attachment.Price)

	return e.Bytes()
}

func (tx *AliasSell) SizeInBytes() int {
	return HeaderSize(tx.Header) + 1 + len(tx.Attachment.Name) + 8
}
