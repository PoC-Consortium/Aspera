package attachment

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
)

type Attachment interface {
	ToBytes(uint8) ([]byte, error)
}

type attachmentType struct {
	surtype int
	subtype int
	new     func() Attachment
}

var appendixTypeOfName = map[string]func() Attachment{
	"Message":               func() Attachment { return new(SendMessageAttachment) },
	"EncryptedMessage":      func() Attachment { return new(SendMessageAttachment) },
	"PublicKeyAnnouncement": func() Attachment { return new(DummyAttachment) }, // ToDo
	"EncryptToSelfMessage":  func() Attachment { return new(DummyAttachment) }, // ToDo
}
var typeOfName = map[string]*attachmentType{
	"OrdinaryPayment":            &attachmentType{surtype: 0, subtype: 0, new: func() Attachment { return new(DummyAttachment) }},
	"MultiOutCreation":           &attachmentType{surtype: 0, subtype: 1, new: func() Attachment { return new(SendMoneyMultiAttachment) }},
	"MultiSameOutCreation":       &attachmentType{surtype: 0, subtype: 2, new: func() Attachment { return new(SendMoneyMultiSameAttachment) }},
	"ArbitaryMessage":            &attachmentType{surtype: 1, subtype: 0, new: func() Attachment { return new(SendMessageAttachment) }},
	"AliasAssignment":            &attachmentType{surtype: 1, subtype: 1, new: func() Attachment { return new(SetAliasAttachment) }},
	"AccountInfo":                &attachmentType{surtype: 1, subtype: 5, new: func() Attachment { return new(SetAccountInfoAttachment) }},
	"AliasSell":                  &attachmentType{surtype: 1, subtype: 6, new: func() Attachment { return new(SellAliasAttachment) }},
	"AliasBuy":                   &attachmentType{surtype: 1, subtype: 7, new: func() Attachment { return new(BuyAliasAttachment) }},
	"AssetIssuance":              &attachmentType{surtype: 2, subtype: 0, new: func() Attachment { return new(IssueAssetAttachment) }},
	"AssetTransfer":              &attachmentType{surtype: 2, subtype: 1, new: func() Attachment { return new(TransferAssetAttachment) }},
	"AskOrderPlacement":          &attachmentType{surtype: 2, subtype: 2, new: func() Attachment { return new(PlaceAskOrderAttachment) }},
	"BidOrderPlacement":          &attachmentType{surtype: 2, subtype: 3, new: func() Attachment { return new(PlaceBidOrderAttachment) }},
	"AskOrderCancellation":       &attachmentType{surtype: 2, subtype: 4, new: func() Attachment { return new(CancelAskOrderAttachment) }},
	"BidOrderCancellation":       &attachmentType{surtype: 2, subtype: 5, new: func() Attachment { return new(CancelBidOrderAttachment) }},
	"DigitalGoodsListing":        &attachmentType{surtype: 3, subtype: 0, new: func() Attachment { return new(DgsListingAttachment) }},
	"DigitalGoodsDelisting":      &attachmentType{surtype: 3, subtype: 1, new: func() Attachment { return new(DgsDelistingAttachment) }},
	"DigitalGoodsPriceChange":    &attachmentType{surtype: 3, subtype: 2, new: func() Attachment { return new(DgsPriceChangeAttachment) }},
	"DigitalGoodsQuantityChange": &attachmentType{surtype: 3, subtype: 3, new: func() Attachment { return new(DgsQuantityChangeAttachment) }},
	"DigitalGoodsPurchase":       &attachmentType{surtype: 3, subtype: 4, new: func() Attachment { return new(DgsPurchaseAttachment) }},
	"DigitalGoodsDelivery":       &attachmentType{surtype: 3, subtype: 5, new: func() Attachment { return new(DgsDeliveryAttachment) }},
	"DigitalGoodsFeedback":       &attachmentType{surtype: 3, subtype: 6, new: func() Attachment { return new(DgsFeedbackAttachment) }},
	"DigitalGoodsRefund":         &attachmentType{surtype: 3, subtype: 7, new: func() Attachment { return new(DgsRefundAttachment) }},
	"EffectiveBalanceLeasing":    &attachmentType{surtype: 4, subtype: 0, new: func() Attachment { return new(LeaseBalanceAttachment) }},
	"RewardRecipientAssignment":  &attachmentType{surtype: 20, subtype: 0, new: func() Attachment { return new(SetRewardRecipientAttachment) }},
	"EscrowCreation":             &attachmentType{surtype: 21, subtype: 0, new: func() Attachment { return new(SendMoneyEscrowAttachment) }},
	"EscrowSign":                 &attachmentType{surtype: 21, subtype: 1, new: func() Attachment { return new(EscrowSignAttachment) }},
	//"EscrowResult":                  &attachmentType{surtype: 21, subtype: 2, new: func() Attachment { return new() }},
	"SubscriptionSubscribe": &attachmentType{surtype: 21, subtype: 3, new: func() Attachment { return new(SendMoneySubscriptionAttachment) }},
	"SubscriptionCancel":    &attachmentType{surtype: 21, subtype: 4, new: func() Attachment { return new(SubscriptionCancelAttachment) }},
	//"SubscriptionPayment":           &attachmentType{surtype: 21, subtype: 5, new: func() Attachment { return new() }},
	//"AutomatedTransactionsCreation": &attachmentType{surtype: 22, subtype: 0, new: func() Attachment { return new() }},
	//"AutomatedTransactionsPayment":  &attachmentType{surtype: 22, subtype: 1, new: func() Attachment { return new() }}, // AT Payment
	//"PublicKeyAnnouncement"
	//"EncryptToSelfMessage"
	// Appendix Only Type :-(
	"Message": &attachmentType{new: func() Attachment { return new(SendMessageAttachment) }},
}
var typeFor = make(map[uint16]*attachmentType)

func init() {
	for _, a := range typeOfName {
		typeFor[uint16(a.surtype)<<4|uint16(a.subtype)] = a
	}
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

func GuessAttachmentsAndAppendicesFromJSON(bs []byte) ([]Attachment, error) {
	var err error

	var txJSON *gabs.Container
	if txJSON, err = gabs.ParseJSON(bs); err != nil {
		return nil, err
	}

	if children, err := txJSON.S("attachment").ChildrenMap(); err != nil {
		return nil, err
	} else if len(children) == 0 {
		return []Attachment{new(DummyAttachment)}, nil
	}
	attachmentType, exists := typeFor[uint16(txJSON.Path("type").Data().(float64))<<4|uint16(txJSON.Path("subtype").Data().(float64))]
	if exists {
		attachments := []Attachment{attachmentType.new()}
		for appendixName, f := range appendixTypeOfName {
			appendixIdentifier := "version." + appendixName
			if txJSON.Exists("attachment", appendixIdentifier) {
				attachments = append(attachments, f())
			}
		}
		return attachments, nil
	}

	return nil, errors.New("tx attachment is not implemented for: " + txJSON.String())
}

type UInt64StringSlice []uint64

func (slice UInt64StringSlice) MarshalJSON() ([]byte, error) {
	values := make([]string, len(slice))
	for i, value := range []uint64(slice) {
		values[i] = fmt.Sprintf(`"%v"`, value)
	}

	return []byte(fmt.Sprintf("[%v]", strings.Join(values, ","))), nil
}

func (slice *UInt64StringSlice) UnmarshalJSON(b []byte) error {
	// Try array of strings first.
	var values []string
	err := json.Unmarshal(b, &values)
	if err != nil {
		// Fall back to array of integers:
		var values []uint64
		if err := json.Unmarshal(b, &values); err != nil {
			return err
		}
		*slice = values
		return nil
	}
	*slice = make([]uint64, len(values))
	for i, value := range values {
		value, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		(*slice)[i] = value
	}
	return nil
}
