package attachment

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"

	"github.com/ac0v/aspera/pkg/parsing"
)

type Attachment interface {
	ToBytes(version uint8) ([]byte, error)
	FromBytes(bs []byte, version uint8) (int, error)
}

type attachmentType struct {
	surtype           int
	subtype           int
	new               func() Attachment
	supersedeAppendix string
}

var appendixTypeOfName = map[string]func() Attachment{
	"Message":               func() Attachment { return new(Message) },
	"EncryptedMessage":      func() Attachment { return new(EncryptedMessage) },
	"PublicKeyAnnouncement": func() Attachment { return new(PublicKeyAnnouncement) },
	"EncryptToSelfMessage":  func() Attachment { return new(EncryptedToSelfMessage) },
}

var typeOfName = map[string]*attachmentType{
	"OrdinaryPayment":               &attachmentType{surtype: 0, subtype: 0, new: func() Attachment { return new(Dummy) }},
	"MultiOutCreation":              &attachmentType{surtype: 0, subtype: 1, new: func() Attachment { return new(SendMoneyMulti) }},
	"MultiSameOutCreation":          &attachmentType{surtype: 0, subtype: 2, new: func() Attachment { return new(SendMoneyMultiSame) }},
	"ArbitaryMessage":               &attachmentType{surtype: 1, subtype: 0, new: func() Attachment { return new(Message) }, supersedeAppendix: "Message"},
	"AliasAssignment":               &attachmentType{surtype: 1, subtype: 1, new: func() Attachment { return new(SetAlias) }},
	"AccountInfo":                   &attachmentType{surtype: 1, subtype: 5, new: func() Attachment { return new(SetAccountInfo) }},
	"AliasSell":                     &attachmentType{surtype: 1, subtype: 6, new: func() Attachment { return new(SellAlias) }},
	"AliasBuy":                      &attachmentType{surtype: 1, subtype: 7, new: func() Attachment { return new(BuyAlias) }},
	"AssetIssuance":                 &attachmentType{surtype: 2, subtype: 0, new: func() Attachment { return new(IssueAsset) }},
	"AssetTransfer":                 &attachmentType{surtype: 2, subtype: 1, new: func() Attachment { return new(TransferAsset) }},
	"AskOrderPlacement":             &attachmentType{surtype: 2, subtype: 2, new: func() Attachment { return new(PlaceAskOrder) }},
	"BidOrderPlacement":             &attachmentType{surtype: 2, subtype: 3, new: func() Attachment { return new(PlaceBidOrder) }},
	"AskOrderCancellation":          &attachmentType{surtype: 2, subtype: 4, new: func() Attachment { return new(CancelAskOrder) }},
	"BidOrderCancellation":          &attachmentType{surtype: 2, subtype: 5, new: func() Attachment { return new(CancelBidOrder) }},
	"DigitalGoodsListing":           &attachmentType{surtype: 3, subtype: 0, new: func() Attachment { return new(DgsListing) }},
	"DigitalGoodsDelisting":         &attachmentType{surtype: 3, subtype: 1, new: func() Attachment { return new(DgsDelisting) }},
	"DigitalGoodsPriceChange":       &attachmentType{surtype: 3, subtype: 2, new: func() Attachment { return new(DgsPriceChange) }},
	"DigitalGoodsQuantityChange":    &attachmentType{surtype: 3, subtype: 3, new: func() Attachment { return new(DgsQuantityChange) }},
	"DigitalGoodsPurchase":          &attachmentType{surtype: 3, subtype: 4, new: func() Attachment { return new(DgsPurchase) }},
	"DigitalGoodsDelivery":          &attachmentType{surtype: 3, subtype: 5, new: func() Attachment { return new(DgsDelivery) }},
	"DigitalGoodsFeedback":          &attachmentType{surtype: 3, subtype: 6, new: func() Attachment { return new(DgsFeedback) }},
	"DigitalGoodsRefund":            &attachmentType{surtype: 3, subtype: 7, new: func() Attachment { return new(DgsRefund) }},
	"EffectiveBalanceLeasing":       &attachmentType{surtype: 4, subtype: 0, new: func() Attachment { return new(LeaseBalance) }},
	"RewardRecipientAssignment":     &attachmentType{surtype: 20, subtype: 0, new: func() Attachment { return new(SetRewardRecipient) }},
	"EscrowCreation":                &attachmentType{surtype: 21, subtype: 0, new: func() Attachment { return new(SendMoneyEscrow) }},
	"EscrowSign":                    &attachmentType{surtype: 21, subtype: 1, new: func() Attachment { return new(EscrowSign) }},
	"EscrowResult":                  &attachmentType{surtype: 21, subtype: 2, new: func() Attachment { return new(EscrowResult) }},
	"SubscriptionSubscribe":         &attachmentType{surtype: 21, subtype: 3, new: func() Attachment { return new(SendMoneySubscription) }},
	"SubscriptionCancel":            &attachmentType{surtype: 21, subtype: 4, new: func() Attachment { return new(SubscriptionCancel) }},
	"SubscriptionPayment":           &attachmentType{surtype: 21, subtype: 5, new: func() Attachment { return new(AdvancedPaymentSubscriptionPayment) }},
	"AutomatedTransactionsCreation": &attachmentType{surtype: 22, subtype: 0, new: func() Attachment { return new(AutomatedTransactionsCreation) }},
	"AutomatedTransactionsPayment":  &attachmentType{surtype: 22, subtype: 1, new: func() Attachment { return new(AutomatedTransactionsPayment) }},
}
var typeFor = make(map[uint16]*attachmentType)

func init() {
	for _, a := range typeOfName {
		typeFor[uint16(a.surtype)<<4|uint16(a.subtype)] = a
	}
}

func FromBytes(bs []byte, surtype, subtype, version uint8, flags uint32) ([]Attachment, error) {
	attachmentType, exists := typeFor[uint16(surtype)<<4|uint16(subtype)]
	if !exists {
		return nil, fmt.Errorf("no parse function for transaction with type %d and subtype %d",
			surtype, subtype)
	}
	attachment := attachmentType.new()
	attachmentLen, err := attachment.FromBytes(bs, version)
	if err != nil {
		return nil, err
	}

	attachments := []Attachment{attachment}
	if flags == 0 {
		return attachments, err
	}

	remainingBs := bs[attachmentLen:]
	if flags&(1<<0) != 0 {
		if version > 0 {
			if err := parsing.SkipByteInSlice(&remainingBs); err != nil {
				return nil, err
			}
		}

		message := new(Message)
		len, err := message.FromBytes(remainingBs, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, message)

		remainingBs = remainingBs[:len]
	}

	if flags&(1<<1) != 0 {
		if version > 0 {
			if err := parsing.SkipByteInSlice(&remainingBs); err != nil {
				return nil, err
			}
		}

		encryptedMessage := new(EncryptedMessage)
		len, err := encryptedMessage.FromBytes(remainingBs, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, encryptedMessage)

		remainingBs = remainingBs[:len]
	}

	if flags&(1<<2) != 0 {
		if version > 0 {
			if err := parsing.SkipByteInSlice(&remainingBs); err != nil {
				return nil, err
			}
		}

		publicKeyAnnouncement := new(PublicKeyAnnouncement)
		len, err := publicKeyAnnouncement.FromBytes(remainingBs, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, publicKeyAnnouncement)

		remainingBs = remainingBs[:len]
	}

	if flags&(1<<3) != 0 {
		if version > 0 {
			if err := parsing.SkipByteInSlice(&remainingBs); err != nil {
				return nil, err
			}
		}

		encryptedToSelfMessage := new(EncryptedToSelfMessage)
		len, err := encryptedToSelfMessage.FromBytes(remainingBs, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, encryptedToSelfMessage)

		remainingBs = remainingBs[:len]
	}

	return attachments, nil
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
		return []Attachment{new(Dummy)}, nil
	}
	attachmentType, exists := typeFor[uint16(txJSON.Path("type").Data().(float64))<<4|uint16(txJSON.Path("subtype").Data().(float64))]
	if exists {
		attachments := []Attachment{attachmentType.new()}
		if err := json.Unmarshal(txJSON.S("attachment").Bytes(), attachments[0]); err != nil {
			return nil, err
		}
		for appendixName, f := range appendixTypeOfName {
			appendixIdentifier := "version." + appendixName
			if txJSON.Exists("attachment", appendixIdentifier) && attachmentType.supersedeAppendix != appendixName {
				appendix := f()
				if err := json.Unmarshal(txJSON.S("attachment").Bytes(), appendix); err != nil {
					return nil, err
				}
				attachments = append(attachments, appendix)

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
