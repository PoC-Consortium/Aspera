package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	EffectiveBalanceLeasingType    = 4
	EffectiveBalanceLeasingSubType = 0
)

type EffectiveBalanceLeasing struct {
	*pb.EffectiveBalanceLeasing
}

func (tx *EffectiveBalanceLeasing) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint32(tx.Attachment.Period)
}

func (tx *EffectiveBalanceLeasing) AttachmentSizeInBytes() int {
	return 4
}

func (tx *EffectiveBalanceLeasing) GetType() uint16 {
	return EffectiveBalanceLeasingSubType<<8 | EffectiveBalanceLeasingType
}
