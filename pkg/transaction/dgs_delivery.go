package transaction

import (
	"bytes"
	"encoding/binary"
	"io"
)

type DgsDeliveryTransaction struct {
	Purchase    uint64
	GoodsLength uint32
	GoodsData   []byte
	GoodsNonce  []byte
	DiscountNQT uint64
}

func DgsDeliveryTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx DgsDeliveryTransaction
	if len(bs) < 16 {
		return nil, 0, io.ErrUnexpectedEOF
	}

	encryptedGoodsLenth := binary.LittleEndian.Uint16(bs[8:10])

	r := bytes.NewReader(bs)

	if err := binary.Read(r, binary.LittleEndian, &tx.Purchase); err != nil {
		return nil, 0, err
	}

	if err := binary.Read(r, binary.LittleEndian, &tx.GoodsLength); err != nil {
		return nil, 0, err
	}

	goodsData := make([]byte, encryptedGoodsLenth)
	if err := binary.Read(r, binary.LittleEndian, &goodsData); err != nil {
		return nil, 0, err
	}
	tx.GoodsData = goodsData

	goodsNonce := make([]byte, 32)
	if err := binary.Read(r, binary.LittleEndian, &goodsNonce); err != nil {
		return nil, 0, err
	}
	tx.GoodsNonce = goodsNonce

	if err := binary.Read(r, binary.LittleEndian, &tx.DiscountNQT); err != nil {
		return nil, 0, err
	}

	return &tx, int(r.Size()) - r.Len(), nil
}
