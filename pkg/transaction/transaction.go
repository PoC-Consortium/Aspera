package transaction

import (
	"errors"
	"strings"

	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/crypto"
	"github.com/ac0v/aspera/pkg/encoding"
	env "github.com/ac0v/aspera/pkg/environment"
	. "github.com/ac0v/aspera/pkg/log"

	"github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	"go.uber.org/zap"
)

const (
	signatureOffset = 96
	signatureLen    = 64

	maxTimestampDifference = 15
)

var (
	ErrTransactionSignatureMismatch = errors.New("transaction signature mismatch")
	ErrInvalidTransactionTimestamp  = errors.New("transaction timestamp invalid")
	ErrInvalidTransactionID         = errors.New("transaction id invalid")
	ErrTransactionFeeTooLow         = errors.New("transaction fee too low")
)

type Transaction interface {
	GetType() uint16
	WriteAttachmentBytes(e encoding.Encoder)
	AttachmentSizeInBytes() int
	GetHeader() *pb.TransactionHeader
	GetAppendix() *pb.Appendix
	proto.Message
}

var typeNameToTransaction = map[string]func() Transaction{
	"p2p.OrdinaryPayment":               func() Transaction { return &OrdinaryPayment{new(pb.OrdinaryPayment)} },
	"p2p.AccountInfo":                   func() Transaction { return &AccountInfo{new(pb.AccountInfo)} },
	"p2p.AliasAssignment":               func() Transaction { return &AliasAssignment{new(pb.AliasAssignment)} },
	"p2p.AliasBuy":                      func() Transaction { return &AliasBuy{new(pb.AliasBuy)} },
	"p2p.AliasSell":                     func() Transaction { return &AliasSell{new(pb.AliasSell)} },
	"p2p.ArbitaryMessage":               func() Transaction { return &ArbitaryMessage{new(pb.ArbitaryMessage)} },
	"p2p.AskOrderCancellation":          func() Transaction { return &AskOrderCancellation{new(pb.AskOrderCancellation)} },
	"p2p.AskOrderPlacement":             func() Transaction { return &AskOrderPlacement{new(pb.AskOrderPlacement)} },
	"p2p.AssetIssuance":                 func() Transaction { return &AssetIssuance{new(pb.AssetIssuance)} },
	"p2p.AssetTransfer":                 func() Transaction { return &AssetTransfer{new(pb.AssetTransfer)} },
	"p2p.AutomatedTransactionsCreation": func() Transaction { return &AutomatedTransactionsCreation{new(pb.AutomatedTransactionsCreation)} },
	"p2p.BidOrderCancellation":          func() Transaction { return &BidOrderCancellation{new(pb.BidOrderCancellation)} },
	"p2p.BidOrderPlacement":             func() Transaction { return &BidOrderPlacement{new(pb.BidOrderPlacement)} },
	"p2p.DigitalGoodsDelisting":         func() Transaction { return &DigitalGoodsDelisting{new(pb.DigitalGoodsDelisting)} },
	"p2p.DigitalGoodsDelivery":          func() Transaction { return &DigitalGoodsDelivery{new(pb.DigitalGoodsDelivery)} },
	"p2p.DigitalGoodsFeedback":          func() Transaction { return &DigitalGoodsFeedback{new(pb.DigitalGoodsFeedback)} },
	"p2p.DigitalGoodsListing":           func() Transaction { return &DigitalGoodsListing{new(pb.DigitalGoodsListing)} },
	"p2p.DigitalGoodsPriceChange":       func() Transaction { return &DigitalGoodsPriceChange{new(pb.DigitalGoodsPriceChange)} },
	"p2p.DigitalGoodsPurchase":          func() Transaction { return &DigitalGoodsPurchase{new(pb.DigitalGoodsPurchase)} },
	"p2p.DigitalGoodsQuantityChange":    func() Transaction { return &DigitalGoodsQuantityChange{new(pb.DigitalGoodsQuantityChange)} },
	"p2p.DigitalGoodsRefund":            func() Transaction { return &DigitalGoodsRefund{new(pb.DigitalGoodsRefund)} },
	"p2p.EffectiveBalanceLeasing":       func() Transaction { return &EffectiveBalanceLeasing{new(pb.EffectiveBalanceLeasing)} },
	"p2p.EscrowCreation":                func() Transaction { return &EscrowCreation{new(pb.EscrowCreation)} },
	"p2p.EscrowResult":                  func() Transaction { return &EscrowResult{new(pb.EscrowResult)} },
	"p2p.EscrowSign":                    func() Transaction { return &EscrowSign{new(pb.EscrowSign)} },
	"p2p.MultiOutCreation":              func() Transaction { return &MultiOutCreation{new(pb.MultiOutCreation)} },
	"p2p.MultiSameOutCreation":          func() Transaction { return &MultiSameOutCreation{new(pb.MultiSameOutCreation)} },
	"p2p.RewardRecipientAssignment":     func() Transaction { return &RewardRecipientAssignment{new(pb.RewardRecipientAssignment)} },
	"p2p.SubscriptionCancel":            func() Transaction { return &SubscriptionCancel{new(pb.SubscriptionCancel)} },
	"p2p.SubscriptionPayment":           func() Transaction { return &SubscriptionPayment{new(pb.SubscriptionPayment)} },
	"p2p.SubscriptionSubscribe":         func() Transaction { return &SubscriptionSubscribe{new(pb.SubscriptionSubscribe)} },
}

func AnyToTransaction(a *any.Any) (Transaction, error) {
	typeName := a.TypeUrl
	if slash := strings.LastIndex(typeName, "/"); slash >= 0 {
		typeName = typeName[slash+1:]
	}
	txNew, knownType := typeNameToTransaction[typeName]
	if !knownType {
		Log.Fatal("unkown transaction type name", zap.String("typeName", typeName))
	}
	tx := txNew()
	err := proto.Unmarshal(a.Value, tx)
	return tx, err
}

func ToBytes(tx Transaction) []byte {
	a := tx.GetAppendix()
	h := tx.GetHeader()
	hasAppendix := a != nil
	var flags uint32
	var appendixSize int
	if hasAppendix {
		appendixSize = AppendixSizeInBytes(a, h.Version)
		flags = AppendixFlags(a)
	}

	txType := tx.GetType()
	e := encoding.NewEncoder(HeaderSizeInBytes(h, txType) + tx.AttachmentSizeInBytes() + appendixSize)

	WriteHeader(e, h, flags, txType)
	tx.WriteAttachmentBytes(e)

	if hasAppendix {
		WriteAppendix(e, a, h.Version)
	}

	return e.Bytes()
}

func CalculateID(txBsWithZeroedSignature []byte) uint64 {
	_, txId := crypto.BytesToHashAndID(txBsWithZeroedSignature)
	return txId
}

func GetExpiration(tx Transaction) uint32 {
	h := tx.GetHeader()
	return h.Timestamp + 60*h.Deadline
}

func Validate(tx Transaction, height int32, blockTimestamp, now uint32) error {
	if err := validateTimestamp(tx, blockTimestamp, now); err != nil {
		return err
	}
	if err := validateFee(tx.GetHeader().Fee, height); err != nil {
		return err
	}

	txBsWithZeroedSignature := ToBytes(tx)
	for i := signatureOffset; i < signatureOffset+signatureLen; i++ {
		txBsWithZeroedSignature[i] = 0
	}

	// TODO: cache tx id
	if err := validateID(txBsWithZeroedSignature); err != nil {
		return err
	}
	if err := validateSignature(tx, txBsWithZeroedSignature); err != nil {
		return err
	}

	return nil
}

func validateFee(fee uint64, height int32) error {
	if fee < env.MinimumFee(height) {
		return ErrTransactionFeeTooLow
	}
	return nil
}

func validateID(txBsWithZeroedSig []byte) error {
	if txID := CalculateID((txBsWithZeroedSig)); txID == 0 {
		return ErrInvalidTransactionID
	}
	return nil
}

func validateTimestamp(tx Transaction, blockTimestamp, now uint32) error {
	txTimestamp := tx.GetHeader().Timestamp
	switch {
	case txTimestamp > now+maxTimestampDifference:
		return ErrInvalidTransactionTimestamp
	case txTimestamp > blockTimestamp+maxTimestampDifference:
		return ErrInvalidTransactionTimestamp
	case GetExpiration(tx) < blockTimestamp:
		return ErrInvalidTransactionTimestamp
	}
	return nil
}

func validateSignature(tx Transaction, txBsWithZeroedSig []byte) error {
	if crypto.Verify(tx.GetHeader().Signature, txBsWithZeroedSig, tx.GetHeader().SenderPublicKey, true) {
		return nil
	}
	return ErrTransactionSignatureMismatch
}

func WriteHeader(e encoding.Encoder, h *pb.TransactionHeader, flags uint32, txType uint16) {
	e.WriteUint16(txType | uint16(h.Version)<<12)
	e.WriteUint32(h.Timestamp)
	e.WriteUint16(uint16(h.Deadline))
	e.WriteBytes(h.SenderPublicKey)
	e.WriteUint64(h.Recipient)
	e.WriteUint64(h.Amount)
	e.WriteUint64(h.Fee)
	if len(h.ReferencedTransactionFullHash) == 0 {
		e.WriteZeros(32)
	} else {
		e.WriteBytes(h.ReferencedTransactionFullHash)
	}
	e.WriteBytes(h.Signature)
	if h.Version > 0 {
		e.WriteUint32(flags)
		e.WriteUint32(h.EcBlockHeight)
		e.WriteUint64(h.EcBlockId)

		// TODO: not an ordinary payment or arbitary message
		if txType != 0 && txType != 1 {
			e.WriteUint8(uint8(h.Version))
		}
	}
}

func HeaderSizeInBytes(h *pb.TransactionHeader, txType uint16) int {
	l := 2 + 4 + 2 + 32 + 8 + 8 + 8 + 32 + 64
	if h.Version > 0 {
		l += 4 + 4 + 8

		// TODO: not an ordinary payment or arbitary message
		if txType != 0 && txType != 1 {
			l++
		}
	}
	return l
}

func WriteAppendix(e encoding.Encoder, a *pb.Appendix, version uint32) {
	if a.Message != nil {
		if version > 0 {
			e.WriteUint8(uint8(version))
		}
		e.WriteStringBytesWithInt32Len(a.Message.IsText, a.Message.Content)
	}
	if a.EncryptedMessage != nil {
		if version > 0 {
			e.WriteUint8(uint8(version))
		}
		e.WriteBytesWithInt32Len(a.EncryptedMessage.IsText, a.EncryptedMessage.Data)
		e.WriteBytes(a.EncryptedMessage.Nonce)
	}
	if a.PublicKeyAnnouncement != nil {
		if version > 0 {
			e.WriteUint8(uint8(version))
		}
		e.WriteBytes(a.PublicKeyAnnouncement.PublicKey)
	}
	if a.EncryptToSelfMessage != nil {
		if version > 0 {
			e.WriteUint8(uint8(version))
		}
		e.WriteBytesWithInt32Len(a.EncryptToSelfMessage.IsText, a.EncryptToSelfMessage.Data)
		e.WriteBytes(a.EncryptToSelfMessage.Nonce)
	}
}

func AppendixSizeInBytes(a *pb.Appendix, version uint32) int {
	var l int
	if a.Message != nil {
		if a.Message.IsText {
			l += 4 + len(a.Message.Content)
		} else {
			l += 4 + len(a.Message.Content)/2
		}
		if version > 0 {
			l++
		}
	}
	if a.EncryptedMessage != nil {
		if a.EncryptedMessage.IsText {
			l += 4 + len(a.EncryptedMessage.Data) + len(a.EncryptedMessage.Nonce)
		} else {
			l += 4 + len(a.EncryptedMessage.Data)/2 + len(a.EncryptedMessage.Nonce)
		}
		if version > 0 {
			l++
		}
	}
	if a.PublicKeyAnnouncement != nil {
		l += len(a.PublicKeyAnnouncement.PublicKey)
		if version > 0 {
			l++
		}
	}
	if a.EncryptToSelfMessage != nil {
		if a.EncryptToSelfMessage.IsText {
			l += 4 + len(a.EncryptToSelfMessage.Data) + len(a.EncryptToSelfMessage.Nonce)
		} else {
			l += 4 + len(a.EncryptToSelfMessage.Data)/2 + len(a.EncryptToSelfMessage.Nonce)
		}
		if version > 0 {
			l++
		}
	}
	return l
}

func AppendixFlags(a *pb.Appendix) uint32 {
	var flags uint32
	if a.Message != nil {
		flags |= 1
	}
	if a.EncryptedMessage != nil {
		flags |= 1 << 1
	}
	if a.PublicKeyAnnouncement != nil {
		flags |= 1 << 2
	}
	if a.EncryptToSelfMessage != nil {
		flags |= 1 << 3
	}
	return flags
}
