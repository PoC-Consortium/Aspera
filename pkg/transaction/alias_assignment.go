package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AliasAssignment struct {
	*pb.AliasAssignment
}

func (tx *AliasAssignment) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint8(uint8(len(tx.Attachment.Alias)))
	e.WriteBytes([]byte(tx.Attachment.Alias))
	e.WriteUint8(uint8(len(tx.Attachment.Uri)))
	e.WriteBytes([]byte(tx.Attachment.Uri))

	return e.Bytes()
}

func (tx *AliasAssignment) SizeInBytes() int {
	return HeaderSize(tx.Header) + 1 + len(tx.Attachment.Alias) + 1 + len(tx.Attachment.Uri)
}
