package parsing

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	maxInt32      = 2147483647
	maxMessageLen = 1000
)

var ErrMessageTooLong = errors.New("message too long")

func SkipByteInReader(r *bytes.Reader) error {
	_, err := r.Seek(1, io.SeekCurrent)
	return err
}

func SkipByteInSlice(bs *[]byte) error {
	if len(*bs) == 0 {
		return io.EOF
	}
	*bs = (*bs)[1:]
	return nil
}

func GetMessageLengthAndType(r *bytes.Reader) (int32, int32, bool, error) {
	var len, isTextAndLen int32
	var isText bool

	if err := binary.Read(r, binary.LittleEndian, &isTextAndLen); err != nil {
		return 0, 0, isText, err
	}

	isText = isTextAndLen < 0
	if isText {
		len = isTextAndLen & maxInt32
	} else {
		len = isTextAndLen
	}

	if len > maxMessageLen {
		return 0, 0, isText, ErrMessageTooLong
	}

	return len, isTextAndLen, isText, nil

}
