package transaction

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
	"gopkg.in/restruct.v1"

	"github.com/ac0v/aspera/pkg/crypto"
	jutils "github.com/ac0v/aspera/pkg/json"
	"github.com/ac0v/aspera/pkg/parsing"
	"github.com/ac0v/aspera/pkg/transaction/attachment"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	pubKeyLen = 64

	signatureLen    = 64
	signatureOffset = 1 + 1 + 4 + 2 + 32 + 8 + 8 + 8 + 32
)

type Transaction struct {
	Header      Header
	Attachments []attachment.Attachment `json:"attachment"`
}

type Header struct {
	Type              uint8 `json:"type"`
	SubtypeAndVersion uint8 `json:"-"`

	Subtype uint8 `struct:"-" json:"subtype"`
	Version uint8 `struct:"-" json:"version"`

	Timestamp       uint32          `json:"timestamp,omitempty"`
	Deadline        uint16          `json:"deadline,omitempty"`
	SenderPublicKey jutils.HexSlice `struct:"[32]uint8" json:"senderPublicKey,omitempty"`
	RecipientID     uint64          `json:"recipient,string,omitempty"`
	AmountNQT       uint64          `json:"amountNQT"`
	FeeNQT          uint64          `json:"feeNQT"`
	// TODO: make pointer
	ReferencedTransactionFullHash jutils.HexSlice `struct:"[32]uint8" json:"referencedTransactionFullHash,omitempty"`
	Signature                     jutils.HexSlice `struct:"[64]uint8" json:"signature,omitempty"`

	Flags         uint32 `struct:"-" json:"-"`
	EcBlockHeight uint32 `struct:"-" json:"ecBlockHeight,omitempty"`
	EcBlockID     uint64 `struct:"-" json:"ecBlockId,string,omitempty"`

	// size of bytes of header
	size int `struct:"-"`
}

type TransactionJSON struct {
	Header
	Attachment interface{} `json:"attachment"`
	//[]attachment.Attachment `json:"attachment,omitempty"`
}

func (tx *Transaction) UnmarshalJSON(bs []byte) error {
	var txJSON TransactionJSON
	var err error

	tx.Attachments, err = attachment.GuessAttachmentsAndAppendicesFromJSON(bs)
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
		// copy all except of the attachment interface related stuff
		if len(src.Field(i).Type().Name()) > 0 {
			srcField := src.Field(i)
			dstField := dst.Field(i)
			dstField.Set(reflect.Value(srcField))
		}
	}
	tx.Header.SetSubtypeAndVersion(txJSON.Subtype, txJSON.Version)

	return nil
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	txJSON := new(TransactionJSON)

	src := reflect.ValueOf(tx).Elem()
	dst := reflect.ValueOf(txJSON).Elem()

	for i := 0; i < src.NumField(); i++ {
		// copy all except of the attachment interface related stuff
		if len(dst.Field(i).Type().Name()) > 0 {
			srcField := src.Field(i)
			dstField := dst.Field(i)
			dstField.Set(reflect.Value(srcField))
		}
	}

	// merge attachment and appendixes (if there are any)
	var gabsJSON = *gabs.New()
	for _, attachment := range tx.Attachments {
		var err error
		var bs []byte
		if bs, err = json.Marshal(attachment); err != nil {
			return nil, err
		}

		var attachmentJSON *gabs.Container
		if attachmentJSON, err = gabs.ParseJSON(bs); err != nil {
			return nil, err
		}
		gabsJSON.Merge(attachmentJSON)
	}

	txJSON.Attachment = gabsJSON.Data()
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
	tmp := h.SubtypeAndVersion
	return (tmp & 0xF0) >> 4
}

func (h *Header) GetSubtype() uint8 {
	tmp := h.SubtypeAndVersion
	return tmp & 0x0F
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

	attachments, err := attachment.FromBytes(
		bs[header.size:], header.Type, header.GetSubtype(), header.GetVersion(), header.Flags)
	if err != nil {
		return nil, err
	}
	tx.Attachments = attachments

	return &tx, nil
}

func (tx *Transaction) ToBytes() ([]byte, error) {
	headerBs, err := tx.Header.ToBytes()
	if err != nil {
		return nil, err
	}

	var attachmentsBs []byte
	for _, attachment := range tx.Attachments {
		bs, err := attachment.ToBytes(tx.Header.GetVersion())
		if err != nil {
			return nil, err
		}
		attachmentsBs = append(attachmentsBs, bs...)
	}

	bs := append(headerBs, attachmentsBs...)

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

		if err := parsing.SkipByteInReader(r); err != nil {
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

func (t *Transaction) VerifySignature() (bool, error) {
	var sig [signatureLen]byte
	if bs, err := t.ToBytes(); err == nil {
		for i := signatureOffset; i < signatureOffset+signatureLen; i++ {
			sig[i-signatureOffset] = bs[i]
			bs[i] = 0
		}
		return crypto.Verify(sig[:], bs, t.Header.SenderPublicKey, true), nil
	} else {
		return false, err
	}

}
