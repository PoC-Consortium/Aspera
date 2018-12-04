package transaction

import (
	pb "github.com/PoC-Consortium/aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/aspera/pkg/encoding"
)

const (
	RewardRecipientAssignmentType    = 20
	RewardRecipientAssignmentSubType = 0
)

type RewardRecipientAssignment struct {
	*pb.RewardRecipientAssignment
}

func (tx *RewardRecipientAssignment) WriteAttachmentBytes(e encoding.Encoder) {}

func (tx *RewardRecipientAssignment) AttachmentSizeInBytes() int {
	return 0
}

func (tx *RewardRecipientAssignment) GetType() uint16 {
	return RewardRecipientAssignmentSubType<<8 | RewardRecipientAssignmentType
}
