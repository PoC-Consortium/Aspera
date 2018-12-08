package transaction

import (
	"math"

	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	AutomatedTransactionsCreationType    = 22
	AutomatedTransactionsCreationSubType = 0

	pageSize          = 256
	maxNameLen        = 30
	maxDescriptionLen = 1000
)

type AutomatedTransactionsCreation struct {
	*pb.AutomatedTransactionsCreation
}

func EmptyAutomatedTransactionCreation() *AutomatedTransactionsCreation {
	return &AutomatedTransactionsCreation{
		AutomatedTransactionsCreation: &pb.AutomatedTransactionsCreation{
			Attachment: &pb.AutomatedTransactionsCreation_Attachment{},
		},
	}
}

func (tx *AutomatedTransactionsCreation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes(tx.Attachment.Name)
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes(tx.Attachment.Description)
	e.WriteBytes(tx.Attachment.Bytes)
}

func (tx *AutomatedTransactionsCreation) ReadAttachmentBytes(d encoding.Decoder) {
	nameLen := d.ReadUint8()
	if nameLen > maxNameLen {
		return
	}
	tx.Attachment.Name = d.ReadBytes(int(nameLen))

	descriptionLen := d.ReadUint16()
	if descriptionLen > maxDescriptionLen {
		return
	}
	tx.Attachment.Description = d.ReadBytes(int(descriptionLen))

	startPos := d.Position()
	d.ReadUint16()
	d.ReadUint16()

	codePages := int(d.ReadUint16())
	dataPages := int(d.ReadUint16())

	d.ReadUint16()
	d.ReadUint16()
	d.ReadUint64()

	var codeLen int
	if codePages*pageSize < pageSize+1 {
		codeLen = int(d.ReadUint8())
	} else if codePages*pageSize < math.MaxUint16 {
		codeLen = int(d.ReadUint16())
	} else {
		codeLen = int(d.ReadUint32())
	}
	d.ReadBytes(codeLen)

	var dataLen int
	if dataPages*pageSize < 257 {
		dataLen = int(d.ReadUint8())
	} else if dataPages*pageSize < math.MaxUint16 {
		dataLen = int(d.ReadUint16())
	} else {
		dataLen = int(d.ReadUint32())
	}
	d.ReadBytes(dataLen)

	creationBytesLen := d.Position() - startPos
	d.Reset(startPos)
	tx.Attachment.Bytes = d.ReadBytes(creationBytesLen)
}

func (tx *AutomatedTransactionsCreation) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description) + len(tx.Attachment.Bytes)
}

func (tx *AutomatedTransactionsCreation) GetType() uint16 {
	return AutomatedTransactionsCreationSubType<<8 | AutomatedTransactionsCreationType
}

func (tx *AutomatedTransactionsCreation) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *AutomatedTransactionsCreation) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
