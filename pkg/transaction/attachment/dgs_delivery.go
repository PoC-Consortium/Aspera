package attachment

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"

	jutils "github.com/ac0v/aspera/pkg/json"
	"github.com/ac0v/aspera/pkg/parsing"
	"gopkg.in/restruct.v1"
)

type DgsDelivery struct {
	Purchase          uint64          `json:"purchase,omitempty,string"`
	GoodsIsTextAndLen int32           `json:"-"`
	GoodsIsText       bool            `struct:"-" json:"goodsIsText"`
	GoodsData         jutils.HexSlice `json:"goodsData,omitempty"`
	GoodsNonce        jutils.HexSlice `json:"goodsNonce,omitempty"`
	DiscountNQT       uint64          `json:"discountNQT"`
	Version           int8            `struct:"-" json:"version.DigitalGoodsDelivery,omitempty"`
}

func (attachment *DgsDelivery) FromBytes(bs []byte, version uint8) (int, error) {
	if len(bs) < 16 {
		return 0, io.ErrUnexpectedEOF
	}

	r := bytes.NewReader(bs)

	if err := binary.Read(r, binary.LittleEndian, &attachment.Purchase); err != nil {
		return 0, err
	}

	len, isTextAndLen, isText, err := parsing.GetMessageLengthAndType(r)
	if err != nil {
		return 0, err
	}
	attachment.GoodsIsTextAndLen = isTextAndLen
	attachment.GoodsIsText = isText

	attachment.GoodsData = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &attachment.GoodsData); err != nil {
		return 0, err
	}

	attachment.GoodsNonce = make([]byte, 32)
	if err := binary.Read(r, binary.LittleEndian, &attachment.GoodsNonce); err != nil {
		return 0, err
	}
	if err := binary.Read(r, binary.LittleEndian, &attachment.DiscountNQT); err != nil {
		return 0, err
	}

	return 8 + 4 + int(len) + 32 + 8, nil
}

func (attachment *DgsDelivery) ToBytes(version uint8) ([]byte, error) {
	attachment.GoodsIsTextAndLen = int32(len(attachment.GoodsData))
	if attachment.GoodsIsText {
		attachment.GoodsIsTextAndLen |= math.MinInt32
	}

	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *DgsDelivery) GetFlag() uint32 {
	return StandardAttachmentFlag
}
