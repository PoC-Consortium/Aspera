package transaction

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

const (
	AssetTransferType    = 2
	AssetTransferSubType = 1
)

type AssetTransfer struct {
	*pb.AssetTransfer
}

func EmptyAssetTransfer() *AssetTransfer {
	return &AssetTransfer{
		AssetTransfer: &pb.AssetTransfer{
			Attachment: &pb.AssetTransfer_Attachment{},
		},
	}
}

func (tx *AssetTransfer) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
	if tx.Header.Version == 0 {
		e.WriteUint16(uint16(len(tx.Attachment.Comment)))
		e.WriteBytes(tx.Attachment.Comment)
	}
}

func (tx *AssetTransfer) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Asset = d.ReadUint64()
	tx.Attachment.Quantity = d.ReadUint64()
	if tx.Header.Version == 0 {
		commentLen := d.ReadUint16()
		tx.Attachment.Comment = d.ReadBytes(int(commentLen))
	}
}

func (tx *AssetTransfer) AttachmentSizeInBytes() int {
	if tx.Header.Version == 0 {
		return 8 + 8 + 2 + len(tx.Attachment.Comment)
	} else {
		return 8 + 8
	}
}

func (tx *AssetTransfer) GetType() uint16 {
	return AssetTransferSubType<<8 | AssetTransferType
}

func (tx *AssetTransfer) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *AssetTransfer) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
