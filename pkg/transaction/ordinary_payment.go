package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	OrdinaryPaymentType    = 0
	OrdinaryPaymentSubType = 0
)

type OrdinaryPayment struct {
	*pb.OrdinaryPayment
}

func EmptyOrdinaryPayment() *OrdinaryPayment {
	return &OrdinaryPayment{
		OrdinaryPayment: &pb.OrdinaryPayment{},
	}
}

func (tx *OrdinaryPayment) WriteAttachmentBytes(e encoding.Encoder) {}

func (tx *OrdinaryPayment) AttachmentSizeInBytes() int {
	return 0
}

func (tx *OrdinaryPayment) ReadAttachmentBytes(d encoding.Decoder) {}

func (tx *OrdinaryPayment) GetType() uint16 {
	return OrdinaryPaymentSubType<<8 | OrdinaryPaymentType
}

func (tx *OrdinaryPayment) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *OrdinaryPayment) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
