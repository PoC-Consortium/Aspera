package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	AccountInfoType    = 1
	AccountInfoSubType = 5
)

type AccountInfo struct {
	*pb.AccountInfo
}

func (tx *AccountInfo) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes(tx.Attachment.Name)
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes(tx.Attachment.Description)
}

func EmptyAccountInfo() *AccountInfo {
	return &AccountInfo{
		AccountInfo: &pb.AccountInfo{
			Attachment: &pb.AccountInfo_Attachment{},
		},
	}
}

func (tx *AccountInfo) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description)
}

func (tx *AccountInfo) GetType() uint16 {
	return AccountInfoSubType<<8 | AccountInfoType
}
