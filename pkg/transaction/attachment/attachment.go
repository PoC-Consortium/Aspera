package attachment

import (
	"fmt"
)

type Attachment interface{}

var attachmentParserOf = map[uint16]func([]byte) (Attachment, int, error){
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

func FromBytes(bs []byte, surType, subType uint8) (Attachment, int, error) {
	parse, exists := attachmentParserOf[uint16(surType)<<4|uint16(subType)]
	if !exists {
		return nil, 0, fmt.Errorf("no parse function for transaction with type %d and subtype %d",
			subType, subType)
	}
	return parse(bs)
}
