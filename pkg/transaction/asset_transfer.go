package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AssetTransfer struct {
	*pb.AssetTransfer
}

func (tx *AssetTransfer) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.Header)

	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)

	return e.Bytes()
}

func (tx *AssetTransfer) SizeInBytes() int {
	return HeaderSize(tx.Header) + 8 + 8
}
