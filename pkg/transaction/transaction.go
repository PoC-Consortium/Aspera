package transaction

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Transaction interface{}

type TransactionHeader struct {
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

var transactionParserOf = map[uint16]func([]byte) (Transaction, error){
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

func TransactionFromBytes(bs []byte) (Transaction, error) {
	header, err := headerFromBytes(bs)
	if err != nil {
		return nil, err
	}

	parse, exists := transactionParserOf[uint16(header.Type)<<4|uint16(header.Subtype)]
	if !exists {
		panic("no parse function for transaction with type " + fmt.Sprint(header.Type) +
			" and subtype " + fmt.Sprint(header.Subtype))
	}

	return parse(bs[header.size:])
}

func headerFromBytes(bs []byte) (*TransactionHeader, error) {
	var header TransactionHeader

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

	// TODO: copied from java code, but makes test fail
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
		if _, err := r.Seek(1, io.SeekCurrent); err != nil {
			return nil, err
		}
	}

	header.size = int(r.Size()) - r.Len()

	return &header, nil
}
