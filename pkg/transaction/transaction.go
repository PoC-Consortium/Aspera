package transaction

import (
	"bytes"
	"encoding/binary"

	"github.com/ac0v/aspera/pkg/parsing"
	"github.com/ac0v/aspera/pkg/transaction/appendicies"
	"github.com/ac0v/aspera/pkg/transaction/attachment"
)

type Transaction struct {
	Header     *Header
	Appendices *appendicies.Appendices
	Attachment attachment.Attachment
}

type Header struct {
	Type                          uint8
	Subtype                       uint8
	Timestamp                     uint32
	Deadline                      uint16
	SenderPublicKey               []byte
	RecipientID                   uint64
	AmountNQT                     uint64
	FeeNQT                        uint64
	ReferencedTransactionFullHash []byte
	Signature                     []byte
	Version                       uint8
	Flags                         uint32
	EcBlockHeight                 uint32
	EcBlockID                     uint64

	// size of bytes of header
	size int
}

func TransactionFromBytes(bs []byte) (*Transaction, error) {
	var tx Transaction

	header, err := headerFromBytes(bs)
	if err != nil {
		return nil, err
	}
	tx.Header = header

	attachment, attachmentLen, err := attachment.FromBytes(bs[header.size:], header.Type, header.Subtype)
	if err != nil {
		return nil, err
	}
	tx.Attachment = attachment

	appendencies, err := appendicies.FromBytes(bs[header.size+attachmentLen:], header.Flags, header.Version)
	if err != nil {
		return nil, err
	}
	tx.Appendices = appendencies

	return &tx, nil
}

func headerFromBytes(bs []byte) (*Header, error) {
	var header Header

	r := bytes.NewReader(bs)

	if err := binary.Read(r, binary.LittleEndian, &header.Type); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &header.Subtype); err != nil {
		return nil, err
	}

	header.Version = (header.Subtype & 0xF0) >> 4
	header.Subtype = header.Subtype & 0x0F

	if err := binary.Read(r, binary.LittleEndian, &header.Timestamp); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &header.Deadline); err != nil {
		return nil, err
	}

	// TODO: for some transactions the buffer containts sender id instead of sender public key
	header.SenderPublicKey = make([]byte, 32)
	if err := binary.Read(r, binary.LittleEndian, &header.SenderPublicKey); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &header.RecipientID); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &header.AmountNQT); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &header.FeeNQT); err != nil {
		return nil, err
	}

	header.ReferencedTransactionFullHash = make([]byte, 32)
	if err := binary.Read(r, binary.LittleEndian, &header.ReferencedTransactionFullHash); err != nil {
		return nil, err
	}

	header.Signature = make([]byte, 64)
	if err := binary.Read(r, binary.LittleEndian, &header.Signature); err != nil {
		return nil, err
	}

	if header.Version > 0 {
		if err := binary.Read(r, binary.LittleEndian, &header.Flags); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.LittleEndian, &header.EcBlockHeight); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.LittleEndian, &header.EcBlockID); err != nil {
			return nil, err
		}

		if err := parsing.SkipByte(r); err != nil {
			return nil, err
		}
	}

	header.size = int(r.Size()) - r.Len()

	return &header, nil
}
