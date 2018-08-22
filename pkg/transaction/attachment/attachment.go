package attachment

import (
	"fmt"
)

type Attachment interface {
	ToBytes(uint8) ([]byte, error)
}

var attachmentParserOf = map[uint16]func([]byte, uint8) (Attachment, int, error){
	0:   SendMoneyAttachmentFromBytes,
	1:   SendMoneyMultiAttachmentFromBytes,
	2:   SendMoneyMultiSameAttachmentFromBytes,
	16:  SendMessageAttachmentFromBytes,
	17:  SetAliasAttachmentFromBytes,
	21:  SetAccountInfoAttachmentFromBytes,
	22:  SellAliasAttachmentFromBytes,
	23:  BuyAliasAttachmentFromBytes,
	32:  IssueAssetAttachmentFromBytes,
	33:  TransferAssetAttachmentFromBytes,
	34:  PlaceAskOrderAttachmentFromBytes,
	35:  PlaceBidOrderAttachmentFromBytes,
	36:  CancelAskOrderAttachmentFromBytes,
	37:  CancelBidOrderAttachmentFromBytes,
	48:  DgsListingAttachmentFromBytes,
	49:  DgsDelistingAttachmentFromBytes,
	50:  DgsPriceChangeAttachmentFromBytes,
	51:  DgsQuantityChangeAttachmentFromBytes,
	52:  DgsPurchaseAttachmentFromBytes,
	53:  DgsDeliveryAttachmentFromBytes,
	54:  DgsFeedbackAttachmentFromBytes,
	55:  DgsRefundAttachmentFromBytes,
	64:  LeaseBalanceAttachmentFromBytes,
	320: SetRewardRecipientAttachmentFromBytes,
	336: SendMoneyEscrowAttachmentFromBytes,
	337: EscrowSignAttachmentFromBytes,
	338: EscrowResultAttachmentFromBytes,
	339: SendMoneySubscriptionAttachmentFromBytes,
	340: SubscriptionCancelAttachmentFromBytes,
	352: AtPaymentAttachmentFromBytes,
}

func FromBytes(bs []byte, surtype, subtype, version uint8) (Attachment, int, error) {
	parse, exists := attachmentParserOf[uint16(surtype)<<4|uint16(subtype)]
	if !exists {
		return nil, 0, fmt.Errorf("no parse function for transaction with type %d and subtype %d",
			surtype, subtype)
	}
	return parse(bs, version)
}
