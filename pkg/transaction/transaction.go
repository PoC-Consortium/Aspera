package transaction

import (
	"bytes"
	"encoding/binary"
	"github.com/ac0v/aspera/pkg/parsing"
	"github.com/ac0v/aspera/pkg/transaction/appendicies"
	"github.com/ac0v/aspera/pkg/transaction/attachment"
	"github.com/json-iterator/go"
	"gopkg.in/restruct.v1"
	"io"
	"reflect"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Transaction struct {
	Header     Header
	Appendices *appendicies.Appendices
	Attachment attachment.Attachment
}

type Header struct {
	Type              uint8 `json:"type"`
	SubtypeAndVersion uint8 `json:"-"`

	Subtype uint8 `struct:"-" json:"subtype"`
	Version uint8 `struct:"-" json:"version"`

	Timestamp                     uint32 `json:"timestamp,omitempty"`
	Deadline                      uint16 `json:"deadline,omitempty"`
	SenderPublicKey               []byte `struct:"[32]uint8" json:"senderPublicKey,omitempty"`
	RecipientID                   uint64 `json:"recipient,string,omitempty"`
	AmountNQT                     uint64 `json:"amountNQT"`
	FeeNQT                        uint64 `json:"feeNQT"`
	ReferencedTransactionFullHash []byte `struct:"[32]uint8" json:"referencedTransactionFullHash,omitempty"`
	Signature                     []byte `struct:"[64]uint8" json:"signature,omitempty"`

	Flags         uint32 `struct:"-" json:"-"`
	EcBlockHeight uint32 `struct:"-" json:"ecBlockHeight"`
	EcBlockID     uint64 `struct:"-" json:"ecBlockId,string"`

	// size of bytes of header
	size int `struct:"-"`
}

type TransactionJSON struct {
	Header
	Appendices *appendicies.Appendices `json:"-"`
	Attachment attachment.Attachment   `json:"attachment,omitempty"`
}

func (tx *Transaction) UnmarshalJSON(bs []byte) error {
	var txJSON TransactionJSON
	var err error

	attachments, err := attachment.GuessAttachmentsAndAppendicesFromJSON(bs)
	txJSON.Attachment = attachments[0]
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, &txJSON)
	if err != nil {
		return err
	}

	src := reflect.ValueOf(&txJSON).Elem()
	dst := reflect.ValueOf(tx).Elem()

	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)
		dstField.Set(reflect.Value(srcField))
	}

	tx.Header.SetSubtypeAndVersion(txJSON.Subtype, txJSON.Version)

	return nil
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	txJSON := new(TransactionJSON)

	src := reflect.ValueOf(tx).Elem()
	dst := reflect.ValueOf(txJSON).Elem()

	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)
		dstField.Set(reflect.Value(srcField))
	}

	txJSON.Version = txJSON.GetVersion()
	txJSON.Subtype = txJSON.GetSubtype()

	// THX for an inconsistent interface - JAVA...
	if tx.Header.RecipientID == 0 && txJSON.Type == 1 && txJSON.Subtype == 6 {
		j, err := json.Marshal(txJSON)
		return append(j[0:len(j)-1], `, "recipient": "0" }`...), err
	}

	return json.Marshal(txJSON)
}

func (h *Header) GetVersion() uint8 {
	return (h.SubtypeAndVersion & 0xF0) >> 4
}

func (h *Header) GetSubtype() uint8 {
	return h.SubtypeAndVersion & 0x0F
}

func (h *Header) SetSubtypeAndVersion(subtype uint8, version uint8) {
	h.SubtypeAndVersion = (version << 4) | subtype
}

func FromBytes(bs []byte) (*Transaction, error) {
	var tx Transaction

	header, err := headerFromBytes(bs)
	if err != nil {
		return nil, err
	}
	tx.Header = *header

	attachment, attachmentLen, err := attachment.FromBytes(bs[header.size:],
		header.Type, header.GetSubtype(), header.GetVersion())
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

	attachmentBs, err := tx.Attachment.ToBytes(tx.Header.GetVersion())
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
