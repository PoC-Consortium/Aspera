package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type DgsListing struct {
	NumName        uint16 `struct:"uint16,sizeof=Name" json:"-"`
	Name           string `json:"name,omitempty"`
	NumDescription uint16 `struct:"uint16,sizeof=Description" json:"-"`
	Description    string `json:"description,omitempty"`
	NumTags        uint16 `struct:"uint16,sizeof=Tags" json:"-"`
	Tags           string `json:"tags"`
	Quantity       uint32 `json:"quantity,omitempty"`
	PriceNQT       uint64 `json:"priceNQT,omitempty"`
	Version        int8   `struct:"-" json:"version.DigitalGoodsListing,omitempty"`
}

func (attachment *DgsListing) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 2 + len(attachment.Name) + 2 + len(attachment.Description) + 2 + len(attachment.Tags) + 4 + 8, err
}

func (attachment *DgsListing) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *DgsListing) GetFlag() uint32 {
	return StandardAttachmentFlag
}
