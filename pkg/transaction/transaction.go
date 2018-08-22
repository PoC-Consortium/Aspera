package transaction

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/ac0v/aspera/pkg/parsing"
	"github.com/ac0v/aspera/pkg/transaction/appendicies"
	"github.com/ac0v/aspera/pkg/transaction/attachment"
	"gopkg.in/restruct.v1"
)

type Transaction struct {
	Header     *Header
	Appendices *appendicies.Appendices
	Attachment attachment.Attachment
}

type Header struct {
	Type                          uint8
	SubtypeAndVersion             uint8
	Timestamp                     uint32
	Deadline                      uint16
	SenderPublicKey               []byte `struct:"[32]uint8"`
	RecipientID                   uint64
	AmountNQT                     uint64
	FeeNQT                        uint64
	ReferencedTransactionFullHash []byte `struct:"[32]uint8"`
	Signature                     []byte `struct:"[64]uint8"`

	Flags         uint32 `struct:"-"`
	EcBlockHeight uint32 `struct:"-"`
	EcBlockID     uint64 `struct:"-"`

	// size of bytes of header
	size int `struct:"-"`
}

func (h *Header) GetVersion() uint8 {
	return (h.SubtypeAndVersion & 0xF0) >> 4
}

func (h *Header) GetSubtype() uint8 {
	return h.SubtypeAndVersion & 0x0F
}

func FromBytes(bs []byte) (*Transaction, error) {
	var tx Transaction

	header, err := headerFromBytes(bs)
	if err != nil {
		return nil, err
	}
	tx.Header = header

	attachment, attachmentLen, err := attachment.FromBytes(bs[header.size:], header.Type, header.GetSubtype())
	if err != nil {
		return nil, err
	}
	tx.Attachment = attachment

	appendencies, err := appendicies.FromBytes(bs[header.size+attachmentLen:], header.Flags, header.GetVersion())
	if err != nil {
		return nil, err
	}
	tx.Appendices = appendencies

	return &tx, nil
}

func (tx *Transaction) ToBytes() ([]byte, error) {
	headerBs, err := tx.Header.ToBytes()
	if err != nil {
		return nil, err
	}

	attachmentBs, err := tx.Attachment.ToBytes()
	if err != nil {
		return nil, err
	}

	appendiciesBs, err := tx.Appendices.ToBytes(tx.Header.GetVersion())
	if err != nil {
		return nil, err
	}

	bs := append(headerBs, attachmentBs...)
	bs = append(bs, appendiciesBs...)

	return bs, nil
}

func headerFromBytes(bs []byte) (*Header, error) {
	var header Header
	if err := restruct.Unpack(bs, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	header.size = 1 + 1 + 4 + 2 + 32 + 8 + 8 + 8 + 32 + 64

	if header.GetVersion() > 0 {
		additionalSize := 4 + 4 + 8 + 1
		if len(bs) < header.size+additionalSize {
			return nil, io.ErrUnexpectedEOF
		}

		r := bytes.NewReader(bs[header.size:])
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

		header.size += additionalSize
	}

	return &header, nil
}

func (h *Header) ToBytes() ([]byte, error) {
	bs, err := restruct.Pack(binary.LittleEndian, h)
	if err != nil {
		return nil, err
	}

	version := h.GetVersion()
	if version > 0 {
		buf := bytes.NewBuffer(nil)

		if err := binary.Write(buf, binary.LittleEndian, h.Flags); err != nil {
			return nil, err
		}

		if err := binary.Write(buf, binary.LittleEndian, h.EcBlockHeight); err != nil {
			return nil, err
		}

		if err := binary.Write(buf, binary.LittleEndian, h.EcBlockID); err != nil {
			return nil, err
		}

		if err := binary.Write(buf, binary.LittleEndian, version); err != nil {
			return nil, err
		}

		return append(bs, buf.Bytes()...), nil
	}

	return bs, nil
}
