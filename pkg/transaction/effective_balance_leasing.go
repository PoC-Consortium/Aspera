package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	EffectiveBalanceLeasingType    = 4
	EffectiveBalanceLeasingSubType = 0
)

type EffectiveBalanceLeasing struct {
	*pb.EffectiveBalanceLeasing
}

func EmptyEffectiveBalanceLeasing() *EffectiveBalanceLeasing {
	return &EffectiveBalanceLeasing{
		EffectiveBalanceLeasing: &pb.EffectiveBalanceLeasing{
			Attachment: &pb.EffectiveBalanceLeasing_Attachment{},
		},
	}
}

func (tx *EffectiveBalanceLeasing) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint16(uint16(tx.Attachment.Period))
}

func (tx *EffectiveBalanceLeasing) AttachmentSizeInBytes() int {
	return 2
}

func (tx *EffectiveBalanceLeasing) GetType() uint16 {
	return EffectiveBalanceLeasingSubType<<8 | EffectiveBalanceLeasingType
}

func (tx *EffectiveBalanceLeasing) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *EffectiveBalanceLeasing) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
