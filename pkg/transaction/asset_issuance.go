package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type AssetIssuence struct {
	*pb.AssetIssuence
}

func (tx *AssetIssuence) ToBytes() []byte {
	e := encoding.NewEncoder([]byte{})

	WriteHeader(e, tx.TransactionHeader)

	e.WriteUint64(tx.Attachment.AssetTransfer)
	e.WriteUint64(tx.Attachment.Qantity)
	e.WriteUint64(tx.Attachment.Price)
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes([]byte(tx.Attachment.Description))
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint8(tx.Attachment.Decimals)
}

func (tx *AssetIssuence) SizeInBytes() int {
	return HeaderSize(tx.TransactionHeader) + 8 + 8 + 8 + 1 + len(tx.Attachment.Name) + 2 +
		len(tx.Attachment.Description) + 8 + 1
}
