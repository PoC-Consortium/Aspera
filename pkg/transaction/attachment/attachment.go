package attachment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"

	"github.com/ac0v/aspera/pkg/parsing"
)

type Attachment interface {
	ToBytes(uint8) ([]byte, error)
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

func FromBytes(bs []byte, surtype, subtype, version uint8, flags uint32) ([]Attachment, error) {
	parse, exists := attachmentParserOf[uint16(surtype)<<4|uint16(subtype)]

	if !exists {
		return nil, fmt.Errorf("no parse function for transaction with type %d and subtype %d",
			surtype, subtype)
	}

	attachment, attachmentLen, err := parse(bs, version)
	if err != nil {
		return nil, err
	}

	attachments := []Attachment{attachment}
	if flags == 0 {
		return attachments, err
	}

	r := bytes.NewReader(bs[attachmentLen:])
	if flags&(1<<0) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}

		message, err := MessageFromBytes(r, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, message)
	}

	if flags&(1<<1) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}

		encryptedMessage, err := EncryptedMessageFromBytes(r, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, encryptedMessage)
	}

	if flags&(1<<2) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}

		publicKeyAnnouncement, err := PublicKeyAnnouncementFromBytes(r, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, publicKeyAnnouncement)
	}

	if flags&(1<<3) != 0 {
		if version > 0 {
			if err := parsing.SkipByte(r); err != nil {
				return nil, err
			}
		}
		encryptedToSelfMessage, err := EncryptedToSelfMessageFromBytes(r, version)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, encryptedToSelfMessage)
	}

	return attachments, nil
}

func ChooseAttachmentFromJSON(bs []byte) (Attachment, error) {
	txJSON, _ := gabs.ParseJSON(bs)
	children, _ := txJSON.S("attachment").ChildrenMap()

	if len(children) == 0 {
		return new(DummyAttachment), nil
	}

	if txJSON.Path("version").Data().(float64) == 0 {
		txType, _ := txJSON.Path("type").Data().(float64)
		txSubtype, _ := txJSON.Path("subtype").Data().(float64)
		switch txType {
		case 0:
			switch txSubtype {
			case 1:
				return new(SendMoneyMultiAttachment), nil
			case 2:
				return new(SendMoneyMultiSameAttachment), nil
			default:
				goto UNKNOWN_ATTACHMENT
			}
		case 1:
			switch txSubtype {
			case 0:
				return new(SendMessageAttachment), nil
			case 1:
				return new(SetAliasAttachment), nil
			case 5:
				return new(SetAccountInfoAttachment), nil
			case 6:
				return new(SellAliasAttachment), nil
			case 7:
				return new(BuyAliasAttachment), nil
			default:
				goto UNKNOWN_ATTACHMENT
			}
		case 2:
			switch txSubtype {
			case 0:
				return new(IssueAssetAttachment), nil
			case 1:
				return new(TransferAssetAttachment), nil
			case 2:
				return new(PlaceAskOrderAttachment), nil
			case 3:
				return new(PlaceBidOrderAttachment), nil
			case 4:
				return new(CancelAskOrderAttachment), nil
			case 5:
				return new(CancelBidOrderAttachment), nil
			default:
				goto UNKNOWN_ATTACHMENT
			}
		case 3:
			switch txSubtype {
			case 0:
				return new(DgsListingAttachment), nil
			case 1:
				return new(DgsDelistingAttachment), nil
			case 2:
				return new(DgsPriceChangeAttachment), nil
			case 3:
				return new(DgsQuantityChangeAttachment), nil
			case 4:
				return new(DgsPurchaseAttachment), nil
			case 5:
				return new(DgsDeliveryAttachment), nil
			case 6:
				return new(DgsFeedbackAttachment), nil
			case 7:
				return new(DgsRefundAttachment), nil
			default:
				goto UNKNOWN_ATTACHMENT
			}
		case 4:
			switch txSubtype {
			case 0:
				return new(LeaseBalanceAttachment), nil
			default:
				goto UNKNOWN_ATTACHMENT
			}
		case 20:
			switch txSubtype {
			case 0:
				return new(SetRewardRecipientAttachment), nil
			default:
				goto UNKNOWN_ATTACHMENT
			}
		case 21:
			switch txSubtype {
			case 0:
				return new(SendMoneyEscrowAttachment), nil
			case 1:
				return new(EscrowSignAttachment), nil
			case 3:
				return new(SendMoneySubscriptionAttachment), nil
			case 4:
				return new(SubscriptionCancelAttachment), nil
			default:
				goto UNKNOWN_ATTACHMENT
			}
		}
	UNKNOWN_ATTACHMENT:
		return nil, errors.New("tx attachment is not implemented for: " + txJSON.String())

		/*
		   private static final byte TYPE_AUTOMATED_TRANSACTIONS = 22;

		   private static final byte SUBTYPE_AT_CREATION = 0;
		   private static final byte SUBTYPE_AT_NXT_PAYMENT = 1;

		   private static final byte SUBTYPE_ACCOUNT_CONTROL_EFFECTIVE_BALANCE_LEASING = 0;

		   private static final byte SUBTYPE_BURST_MINING_REWARD_RECIPIENT_ASSIGNMENT = 0;

		   private static final byte SUBTYPE_ADVANCED_PAYMENT_ESCROW_CREATION = 0;
		   private static final byte SUBTYPE_ADVANCED_PAYMENT_ESCROW_SIGN = 1;
		   private static final byte SUBTYPE_ADVANCED_PAYMENT_ESCROW_RESULT = 2;
		   private static final byte SUBTYPE_ADVANCED_PAYMENT_SUBSCRIPTION_SUBSCRIBE = 3;
		   private static final byte SUBTYPE_ADVANCED_PAYMENT_SUBSCRIPTION_CANCEL = 4;
		   private static final byte SUBTYPE_ADVANCED_PAYMENT_SUBSCRIPTION_PAYMENT = 5;
		*/
	}

	var attachmentType string
	for key, _ := range children {
		if strings.HasPrefix(key, "version.") {
			attachmentType = strings.TrimPrefix(key, "version.")
			break
		}
	}

	switch attachmentType {
	case "Message": // ArbitaryMessage ??
		return new(SendMessageAttachment), nil
	case "MultiOutCreation":
		return new(SendMoneyMultiAttachment), nil
	case "MultiSameOutCreation":
		return new(SendMoneyMultiSameAttachment), nil
	case "AliasAssignment":
		return new(SetAliasAttachment), nil
	case "AliasSell":
		return new(SellAliasAttachment), nil
	case "AliasBuy":
		return new(BuyAliasAttachment), nil
	case "AccountInfo":
		return new(SetAccountInfoAttachment), nil
	case "AssetIssuance":
		return new(IssueAssetAttachment), nil
	case "AssetTransfer":
		return new(TransferAssetAttachment), nil
	case "AskOrderPlacement":
		return new(PlaceAskOrderAttachment), nil
	case "BidOrderPlacement":
		return new(PlaceBidOrderAttachment), nil
	case "AskOrderCancellation":
		return new(CancelAskOrderAttachment), nil
	case "BidOrderCancellation":
		return new(CancelBidOrderAttachment), nil
	case "DigitalGoodsListing":
		return new(DgsListingAttachment), nil
	case "DigitalGoodsDelisting":
		return new(DgsDelistingAttachment), nil
	case "DigitalGoodsPriceChange":
		return new(DgsPriceChangeAttachment), nil
	case "DigitalGoodsQuantityChange":
		return new(DgsQuantityChangeAttachment), nil
	case "DigitalGoodsPurchase":
		return new(DgsPurchaseAttachment), nil
	case "DigitalGoodsDelivery":
		return new(DgsDeliveryAttachment), nil
	case "DigitalGoodsFeedback":
		return new(DgsFeedbackAttachment), nil
	case "DigitalGoodsRefund":
		return new(DgsRefundAttachment), nil
	case "RewardRecipientAssignment":
		return new(SetRewardRecipientAttachment), nil
	case "EscrowCreation":
		return new(SendMoneyEscrowAttachment), nil
	case "EscrowSign":
		return new(EscrowSignAttachment), nil
	case "EscrowResult":
		return new(EscrowResultAttachment), nil
	case "SubscriptionSubscribe":
		return new(SendMoneySubscriptionAttachment), nil
	case "SubscriptionCancel":
		return new(SubscriptionCancelAttachment), nil
	case "AutomatedTransactionsCreation":
		return new(AtPaymentAttachment), nil
	case "EncryptedMessage":
		return new(DummyAttachment), nil // ToDo
	case "EffectiveBalanceLeasing":
		return new(LeaseBalanceAttachment), nil
	case "SubscriptionPayment":
	case "OrdinaryPayment":
	}
	return nil, errors.New(attachmentType + " is not implemented for: " + txJSON.String())
}
