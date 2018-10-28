package transaction

import (
	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
)

const (
	AliasAssignmentType    = 1
	AliasAssignmentSubType = 1
)

type AliasAssignment struct {
	*pb.AliasAssignment
}

func (tx *AliasAssignment) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Alias)))
	e.WriteBytes([]byte(tx.Attachment.Alias))
	e.WriteUint8(uint8(len(tx.Attachment.Uri)))
	e.WriteBytes([]byte(tx.Attachment.Uri))
}

func (tx *AliasAssignment) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Alias) + 1 + len(tx.Attachment.Uri)
}

func (tx *AliasAssignment) GetType() uint16 {
	return AliasAssignmentSubType<<4 | AliasAssignmentType
}
