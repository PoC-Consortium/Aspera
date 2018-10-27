package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AssetIssuence struct {
	*pb.AssetIssuance
}

func (tx *AssetIssuence) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes([]byte(tx.Attachment.Description))
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint8(uint8(tx.Attachment.Decimals))

	return e.Bytes()
}

func (tx *AssetIssuence) SizeInBytes() int {
	return HeaderSize(tx.Header) + 1 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description) + 8 + 1
}
