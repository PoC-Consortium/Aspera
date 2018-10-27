package transaction

import (
	"strings"

	pb "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/encoding"
	. "github.com/ac0v/aspera/pkg/log"

	"github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	"go.uber.org/zap"
)

type Transaction interface {
	ToBytes() []byte
	SizeInBytes() int
	proto.Message
}

var typeNameToTransaction = map[string]func() Transaction{
	"p2p.OrdinaryPayment":               func() Transaction { return &OrdinaryPayment{new(pb.OrdinaryPayment)} },
	"p2p.AccountInfo":                   func() Transaction { return &AliasBuy{new(pb.AliasBuy)} },
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

func WriteHeader(e encoding.Encoder, h *pb.TransactionHeader) {
	e.WriteUint32(h.Timestamp)
	e.WriteUint16(uint16(h.Deadline))
	e.WriteBytes(h.SenderPublicKey)
	e.WriteUint64(h.Recipient)
	e.WriteUint64(h.Amount)
	e.WriteUint64(h.Fee)
	e.WriteBytes(h.ReferencedTransactionFullHash)
	e.WriteBytes(h.Signature)
	if h.Version > 0 {
		// TODO: calc flags
		// e.WriteUint32(h.Flags)
		e.WriteUint32(h.EcBlockHeight)
		e.WriteUint64(h.EcBlockId)
	}
}

func HeaderSize(h *pb.TransactionHeader) int {
	l := 4 + 2 + 32 + 8 + 8 + 8 + 64 + 32 + 4 + 4 + 8
	if h.Version > 0 {
		l += 4 + 4 + 8
	}
	return l
}
