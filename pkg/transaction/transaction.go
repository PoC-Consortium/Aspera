package transaction

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	maxInt32      = 2147483647
	maxMessageLen = 1000
)

var errMessageTooLong = errors.New("message too long")

type Transaction struct {
	Header     *Header
	Appendices *Appendices
	Attachment Attachment
}

type Attachment interface{}

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

type Message struct {
	IsText  bool
	Len     int32
	Content []byte
}

type EncryptedMessage struct {
	IsText bool
	Len    int32
	Data   []byte
	Nonce  []byte
}

type PublicKeyAnnouncement struct {
	PublicKey []byte
}

type EncryptedToSelfMessage struct {
	IsText bool
	Len    int32
	Data   []byte
	Nonce  []byte
}

type Appendices struct {
	Message                *Message
	EncryptedMessage       *EncryptedMessage
	PublicKeyAnnouncement  *PublicKeyAnnouncement
	EncryptedToSelfMessage *EncryptedMessage
}

var transactionParserOf = map[uint16]func([]byte) (Attachment, int, error){
	0:   SendMoneyTransactionFromBytes,
	1:   SendMoneyMultiTransactionFromBytes,
	2:   SendMoneyMultiSameTransactionFromBytes,
	16:  SendMessageTransactionFromBytes,
	17:  SetAliasTransactionFromBytes,
	21:  SetAccountInfoTransactionFromBytes,
	22:  SellAliasTransactionFromBytes,
	23:  BuyAliasTransactionFromBytes,
	32:  IssueAssetTransactionFromBytes,
	33:  TransferAssetTransactionFromBytes,
	34:  PlaceAskOrderTransactionFromBytes,
	35:  PlaceBidOrderTransactionFromBytes,
	36:  CancelAskOrderTransactionFromBytes,
	37:  CancelBidOrderTransactionFromBytes,
	48:  DgsListingTransactionFromBytes,
	49:  DgsDelistingTransactionFromBytes,
	50:  DgsPriceChangeTransactionFromBytes,
	51:  DgsQuantityChangeTransactionFromBytes,
	52:  DgsPurchaseTransactionFromBytes,
	53:  DgsDeliveryTransactionFromBytes,
	54:  DgsFeedbackTransactionFromBytes,
	55:  DgsRefundTransactionFromBytes,
	64:  LeaseBalanceTransactionFromBytes,
	320: SetRewardRecipientTransactionFromBytes,
	336: SendMoneyEscrowTransactionFromBytes,
	337: EscrowSignTransactionFromBytes,
	338: EscrowResultTransactionFromBytes,
	339: SendMoneySubscriptionTransactionFromBytes,
	340: SubscriptionCancelTransactionFromBytes,
	352: AtPaymentTransactionFromBytes,
}

func TransactionFromBytes(bs []byte) (*Transaction, error) {
	var tx Transaction

	header, err := headerFromBytes(bs)
	if err != nil {
		return nil, err
	}
	tx.Header = header

	parse, exists := transactionParserOf[uint16(header.Type)<<4|uint16(header.Subtype)]
	if !exists {
		return nil, fmt.Errorf("no parse function for transaction with type %d and subtype %d",
			header.Type, header.Subtype)
	}

	attachment, attachmentLen, err := parse(bs[header.size:])
	if err != nil {
		return nil, err
	}
	tx.Attachment = attachment

	appendencies, err := appendiciesFromBytes(bs[header.size+attachmentLen:], header.Flags, header.Version)
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

		// skip one byte
		if err := skipByte(r); err != nil {
			return nil, err
		}
	}

	header.size = int(r.Size()) - r.Len()

	return &header, nil
}

func skipByte(r *bytes.Reader) error {
	_, err := r.Seek(1, io.SeekCurrent)
	return err
}

func getMessageLengthAndType(r io.Reader) (int32, bool, error) {
	var len int32
	var isText bool

	if err := binary.Read(r, binary.LittleEndian, &len); err != nil {
		return len, isText, err
	}

	isText = len < 0
	if isText {
		len &= maxInt32
	}

	if len > maxMessageLen {
		return len, isText, errMessageTooLong
	}

	return len, isText, nil
}

func messageFromBytes(r io.Reader) (*Message, error) {
	var message Message

	len, isText, err := getMessageLengthAndType(r)
	if err != nil {
		return nil, err
	}

	message.Len = len
	message.IsText = isText

	message.Content = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &message.Content); err != nil {
		return nil, err
	}

	return &message, nil
}

func encryptedMessageFromBytes(r io.Reader) (*EncryptedMessage, error) {
	var message EncryptedMessage

	len, isText, err := getMessageLengthAndType(r)
	if err != nil {
		return nil, err
	}

	message.Len = len
	message.IsText = isText

	message.Data = make([]byte, len)
	if err := binary.Read(r, binary.LittleEndian, &message.Data); err != nil {
		return nil, err
	}

	message.Nonce = make([]byte, 32)
	err = binary.Read(r, binary.LittleEndian, &message.Nonce)

	return &message, err
}

func publicKeyAnnouncementFromBytes(r io.Reader) (*PublicKeyAnnouncement, error) {
	var message PublicKeyAnnouncement

	message.PublicKey = make([]byte, 32)
	err := binary.Read(r, binary.LittleEndian, &message.PublicKey)

	return &message, err
}

func appendiciesFromBytes(bs []byte, flags uint32, version uint8) (*Appendices, error) {
	var appendicies Appendices

	r := bytes.NewReader(bs)
	if flags&(1<<0) != 0 {
		if version > 0 {
			if err := skipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := messageFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.Message = m
	}

	if flags&(1<<1) != 0 {
		if version > 0 {
			if err := skipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := encryptedMessageFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.EncryptedMessage = m
	}

	if flags&(1<<2) != 0 {
		if version > 0 {
			if err := skipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := publicKeyAnnouncementFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.PublicKeyAnnouncement = m
	}

	if flags&(1<<3) != 0 {
		if version > 0 {
			if err := skipByte(r); err != nil {
				return nil, err
			}
		}
		m, err := encryptedMessageFromBytes(r)
		if err != nil {
			return nil, err
		}
		appendicies.EncryptedToSelfMessage = m
	}

	return &appendicies, nil
}
