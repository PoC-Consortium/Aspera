package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

type Transaction interface {
	ToBytes() []byte
	SizeInBytes() int
}

func WriteHeader(e encoding.Encoder, h *pb.TransactionHeader) {
	e.WriteUint32(h.Timestamp)
	e.WriteUint16(uint16(h.Deadline))
	e.WriteBytes(h.SenderPublicKey)
	e.WriteUint64(h.Recipient)
	e.WriteUint64(h.Amount)
	e.WriteUint64(h.Fee)
	e.WriteBytes(h.ReferencedTransactionFullHash)
	e.WriteBytes(h.Signature)
	if h.Version > 0 {
		// TODO: calc flags
		// e.WriteUint32(h.Flags)
		e.WriteUint32(h.EcBlockHeight)
		e.WriteUint64(h.EcBlockId)
	}
}

func HeaderSize(h *pb.TransactionHeader) int {
	l := 4 + 2 + 32 + 8 + 8 + 8 + 64 + 32 + 4 + 4 + 8
	if h.Version > 0 {
		l += 4 + 4 + 8
	}
	return l
}
