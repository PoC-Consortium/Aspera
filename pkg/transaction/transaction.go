package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type Transaction interface{}

type TransactionHeader struct {
	Type                          uint8  `struct:"uint8"`
	Subtype                       uint8  `struct:"uint8"`
	Timestamp                     uint32 `struct:"uint32"`
	Deadline                      uint16 `struct:"uint16"`
	SenderPublicKey               []byte `struct:"[32]uint8"`
	RecipientID                   uint64 `struct:"uint64"`
	AmountNQT                     uint64 `struct:"uint64"`
	FeeNQT                        uint64 `struct:"uint64"`
	ReferencedTransactionFullHash []byte `struct:"[32]uint8"`
	Signature                     []byte `struct:"[64]uint8"`
	Version                       uint8  `struct:"-"`
}

func headerFromBytes(bs []byte) (*TransactionHeader, error) {
	var header TransactionHeader
	if err := restruct.Unpack(bs, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	// TODO: for some transactions the buffer containts sender id instead of sender public key
	header.Version = (header.Subtype & 0xF0) >> 4
	header.Subtype = header.Subtype & 0x0F

	return &header, nil
}
