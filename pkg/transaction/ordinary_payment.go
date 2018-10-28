package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	OrdinaryPaymentType    = 0
	OrdinaryPaymentSubType = 0
)

type OrdinaryPayment struct {
	*pb.OrdinaryPayment
}

func (tx *OrdinaryPayment) WriteAttachmentBytes(e encoding.Encoder) {}

func (tx *OrdinaryPayment) AttachmentSizeInBytes() int {
	return 0
}

func (tx *OrdinaryPayment) GetType() uint16 {
	return OrdinaryPaymentSubType<<8 | OrdinaryPaymentType
}
