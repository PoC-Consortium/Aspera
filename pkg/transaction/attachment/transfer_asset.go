package attachment

import (
	"bytes"
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type TransferAsset struct {
	Asset       uint64 `json:"asset,omitempty,string"`
	QuantityQNT uint64 `json:"quantityQNT,omitempty"`

	NumComment uint16 `struct:"-" json:"-"`
	Comment    string `struct:"-" json:"comment,omitempty"`
	Version    int8   `struct:"-" json:"version.AssetTransfer,omitempty"`
}

func (attachment *TransferAsset) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)

	len := 8 + 8
	if version == 0 {
		r := bytes.NewReader(bs[len:])
		if err := binary.Read(r, binary.LittleEndian, &attachment.NumComment); err != nil {
			return 0, err
		}

		comment := make([]byte, attachment.NumComment)
		if err := binary.Read(r, binary.LittleEndian, &comment); err != nil {
			return 0, err
		}
		attachment.Comment = string(comment)

		len += 2 + int(attachment.NumComment)
	}

	return len, err
}

func (attachment *TransferAsset) ToBytes(version uint8) ([]byte, error) {
	// TODO: might be better to put this into a unmarshaller...
	attachment.NumComment = uint16(len(attachment.Comment))

	bs, err := restruct.Pack(binary.LittleEndian, attachment)
	if err != nil {
		return nil, err
	}

	if version == 0 {
		buf := bytes.NewBuffer(nil)

		if err := binary.Write(buf, binary.LittleEndian, attachment.NumComment); err != nil {
			return nil, err
		}

		if err := binary.Write(buf, binary.LittleEndian, []byte(attachment.Comment)); err != nil {
			return nil, err
		}

		return append(bs, buf.Bytes()...), nil
	}

	return bs, nil
}

func (attachment *TransferAsset) GetFlag() uint32 {
	return StandardAttachmentFlag
}
