package attachment

import (
	"bytes"
	"encoding/binary"
	"io"

	"gopkg.in/restruct.v1"
)

type DgsDeliveryAttachment struct {
	Purchase    uint64 `json:"purchase,omitempty,string"`
	GoodsLength uint32 `json:"goodsLength,omitempty,string"`
	GoodsIsText bool   `struct:"-" json:"goodsIsText"`
	GoodsData   string `json:"goodsData,omitempty"`
	GoodsNonce  string `json:"goodsNonce,omitempty"`
	DiscountNQT uint64 `json:"discountNQT"`
	Version     int8   `struct:"-" json:"version.DigitalGoodsDelivery,omitempty"`
}

func DgsDeliveryAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment DgsDeliveryAttachment
	if len(bs) < 16 {
		return nil, 0, io.ErrUnexpectedEOF
	}

	encryptedGoodsLenth := binary.LittleEndian.Uint16(bs[8:10])

	r := bytes.NewReader(bs)

	if err := binary.Read(r, binary.LittleEndian, &attachment.Purchase); err != nil {
		return nil, 0, err
	}

	if err := binary.Read(r, binary.LittleEndian, &attachment.GoodsLength); err != nil {
		return nil, 0, err
	}

	goodsData := make([]byte, encryptedGoodsLenth)
	if err := binary.Read(r, binary.LittleEndian, &goodsData); err != nil {
		return nil, 0, err
	}
	attachment.GoodsData = string(goodsData)

	goodsNonce := make([]byte, 32)
	if err := binary.Read(r, binary.LittleEndian, &goodsNonce); err != nil {
		return nil, 0, err
	}
	attachment.GoodsNonce = string(goodsNonce)

	if err := binary.Read(r, binary.LittleEndian, &attachment.DiscountNQT); err != nil {
		return nil, 0, err
	}

	// ToDo:
	// if attachment.GoodsLength < 0 {
	//         attachment.GoodsIsText = true
	// }

	return &attachment, int(r.Size()) - r.Len(), nil
}

func (attachment *DgsDeliveryAttachment) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}
