package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	AliasBuyType    = 1
	AliasBuySubType = 7
)

type AliasBuy struct {
	*pb.AliasBuy
}

func EmptyAliasBuy() *AliasBuy {
	return &AliasBuy{
		AliasBuy: &pb.AliasBuy{
			Attachment: &pb.AliasBuy_Attachment{},
		},
	}
}

func (tx *AliasBuy) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes(tx.Attachment.Name)
}

func (tx *AliasBuy) ReadAttachmentBytes(d encoding.Decoder) {
	nameLen := d.ReadUint8()
	tx.Attachment.Name = d.ReadBytes(int(nameLen))
}

func (tx *AliasBuy) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Name)
}

func (tx *AliasBuy) GetType() uint16 {
	return AliasBuySubType<<8 | AliasBuyType
}

func (tx *AliasBuy) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *AliasBuy) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
