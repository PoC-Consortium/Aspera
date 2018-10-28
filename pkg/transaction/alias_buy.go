package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	AliasBuyType    = 1
	AliasBuySubType = 7
)

type AliasBuy struct {
	*pb.AliasBuy
}

func (tx *AliasBuy) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes([]byte(tx.Attachment.Name))
}

func (tx *AliasBuy) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Name)
}

func (tx *AliasBuy) GetType() uint16 {
	return AliasBuySubType<<8 | AliasBuyType
}
