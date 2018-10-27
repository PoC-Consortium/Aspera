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

	WriteHeader(e, tx.Header)

	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))

	return e.Bytes()
}

func (tx *AliasBuy) SizeInBytes() int {
	return HeaderSize(tx.Header) + 1 + len(tx.Attachment.Name)
}
