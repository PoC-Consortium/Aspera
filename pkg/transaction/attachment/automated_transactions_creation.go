package attachment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	maxAtNameLen       = 30
	maxAtDesriptionLen = 1000

	pageSize = 256
)

type AutomatedTransactionsCreationAttachment struct {
	NumName        uint8  `json:"-"`
	Name           string `json:"name"`
	NumDescription uint16 `json:"-"`
	Description    string `json:"description"`
	CodePages      uint16 `json:"-"`
	DataPages      uint16 `json:"-"`
	Code           string `json:"-"`
	Data           string `json:"-"`
	CreationBytes  string `json:"creationBytes"`
	Version        int8   `struct:"-" json:"version.AutomatedTransactionsCreation,omitempty"`
}

/* TODO: bullshit
After reading Name and Description this function remembers a start position.
It continues reading just to find an end position, throws away everything that was
read and reads everything from start to end into a single buffer.

Either:
1. We count up to end position and still only read one chunk
2. We read all those small buffers and use them later on when executing the at

This needs to be changed as soon as we know which params we need for executing the at.
*/
func AutomatedTransactionsCreationAttachmentFromBytes(bs []byte, version uint8) (Attachment, int, error) {
	var attachment AutomatedTransactionsCreationAttachment

	r := bytes.NewReader(bs)

	if err := binary.Read(r, binary.LittleEndian, &attachment.NumName); err != nil {
		return nil, 0, nil
	}
	if attachment.NumName > maxAtNameLen {
		return nil, 0, fmt.Errorf("at name too long")
	}

	attachmentName := make([]byte, attachment.NumName)
	if err := binary.Read(r, binary.LittleEndian, &attachmentName); err != nil {
		return nil, 0, err
	}
	attachment.Name = string(attachmentName)

	if err := binary.Read(r, binary.LittleEndian, &attachment.NumDescription); err != nil {
		return nil, 0, nil
	}
	if attachment.NumDescription > maxAtDesriptionLen {
		return nil, 0, fmt.Errorf("at description too long")
	}

	attachmentDescription := make([]byte, attachment.NumDescription)
	if err := binary.Read(r, binary.LittleEndian, &attachmentDescription); err != nil {
		return nil, 0, err
	}
	attachment.Description = string(attachmentDescription)

	startPosition := int(r.Size()) - r.Len()

	if _, err := r.Seek(2+2, io.SeekCurrent); err != nil {
		return nil, 0, err
	}

	if err := binary.Read(r, binary.LittleEndian, &attachment.CodePages); err != nil {
		return nil, 0, nil
	}

	if err := binary.Read(r, binary.LittleEndian, &attachment.DataPages); err != nil {
		return nil, 0, nil
	}

	if _, err := r.Seek(2+2+8, io.SeekCurrent); err != nil {
		return nil, 0, err
	}

	var codeLen uint32
	if attachment.CodePages*pageSize < pageSize+1 {
		var codeLenU8 uint8
		if err := binary.Read(r, binary.LittleEndian, &codeLenU8); err != nil {
			return nil, 0, nil
		}
		codeLen = uint32(codeLenU8)
	} else if attachment.CodePages*pageSize < math.MaxInt16+1 {
		var codeLenU16 uint16
		if err := binary.Read(r, binary.LittleEndian, &codeLenU16); err != nil {
			return nil, 0, nil
		}
		codeLen = uint32(codeLenU16)
	} else {
		if err := binary.Read(r, binary.LittleEndian, &codeLen); err != nil {
			return nil, 0, nil
		}
	}

	attachmentCode := make([]byte, codeLen)
	if err := binary.Read(r, binary.LittleEndian, &attachmentCode); err != nil {
		return nil, 0, nil
	}
	attachment.Code = string(attachmentCode)

	var dataLen uint32
	if attachment.DataPages*pageSize < 257 {
		var dataLenU8 uint8
		if err := binary.Read(r, binary.LittleEndian, &dataLenU8); err != nil {
			return nil, 0, nil
		}
		dataLen = uint32(dataLenU8)
	} else if attachment.DataPages*pageSize < math.MaxInt16+1 {
		var dataLenU16 uint16
		if err := binary.Read(r, binary.LittleEndian, &dataLenU16); err != nil {
			return nil, 0, nil
		}
		dataLen = uint32(dataLenU16)
	} else {
		if err := binary.Read(r, binary.LittleEndian, &dataLen); err != nil {
			return nil, 0, nil
		}
	}

	attachmentData := make([]byte, dataLen)
	if err := binary.Read(r, binary.LittleEndian, &attachmentData); err != nil {
		return nil, 0, nil
	}
	attachment.Data = string(attachmentData)

	endPosition := int(r.Size()) - r.Len()
	if _, err := r.Seek(int64(startPosition), io.SeekStart); err != nil {
		return nil, 0, err
	}
	attachmentCreationBytes := make([]byte, endPosition-startPosition)
	if err := binary.Read(r, binary.LittleEndian, &attachmentCreationBytes); err != nil {
		return nil, 0, nil
	}
	attachment.CreationBytes = string(attachmentCreationBytes)

	return &attachment, 0, nil
}

func (attachment *AutomatedTransactionsCreationAttachment) ToBytes(version uint8) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if err := binary.Write(buf, binary.LittleEndian, attachment.NumName); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, []byte(attachment.Name)); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, attachment.NumDescription); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, []byte(attachment.Description)); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, []byte(attachment.CreationBytes)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
