package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	ArbitaryMessageType     = 1
	ArbitaryMessageSubyType = 0
)

type ArbitaryMessage struct {
	*pb.ArbitaryMessage
}

func (tx *ArbitaryMessage) WriteAttachmentBytes(e encoding.Encoder) {}

func (tx *ArbitaryMessage) AttachmentSizeInBytes() int {
	return 0
}

func (tx *ArbitaryMessage) GetType() uint16 {
	return ArbitaryMessageType<<8 | ArbitaryMessageSubyType
}
