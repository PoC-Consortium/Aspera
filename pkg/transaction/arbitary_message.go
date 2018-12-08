package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	ArbitaryMessageType    = 1
	ArbitaryMessageSubType = 0
)

type ArbitaryMessage struct {
	*pb.ArbitaryMessage
}

func EmptyArbitraryMessage() *ArbitaryMessage {
	return &ArbitaryMessage{
		ArbitaryMessage: &pb.ArbitaryMessage{},
	}
}

func (tx *ArbitaryMessage) WriteAttachmentBytes(e encoding.Encoder) {}

func (tx *ArbitaryMessage) AttachmentSizeInBytes() int {
	return 0
}

func (tx *ArbitaryMessage) GetType() uint16 {
	return ArbitaryMessageSubType<<8 | ArbitaryMessageType
}
