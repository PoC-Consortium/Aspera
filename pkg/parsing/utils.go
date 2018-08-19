package parsing

import (
	"bytes"
	"io"
)

func SkipByte(r *bytes.Reader) error {
	_, err := r.Seek(1, io.SeekCurrent)
	return err
}
