package attachment

import ()

type AtPaymentAttachment struct{}

func AtPaymentAttachmentFromBytes(bs []byte) (Attachment, int, error) {
	var attachment AtPaymentAttachment
	return &attachment, 0, nil
}

func (attachment *AtPaymentAttachment) ToBytes() ([]byte, error) {
	return []byte{}, nil
}
