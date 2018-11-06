package encoding

import (
	"encoding/binary"
	"encoding/hex"
	"math"
)

var order = binary.LittleEndian

type Encoder interface {
	WriteUint64(val uint64)
	WriteInt64(val int64)
	WriteUint32(val uint32)
	WriteInt32(val int32)
	WriteUint16(val uint16)
	WriteInt16(val int16)
	WriteUint8(val uint8)
	WriteInt8(val int8)

	WriteBytes(val []byte)

	WriteStringBytesWithInt32Len(isTxt bool, val []byte) error
	WriteBytesWithInt32Len(isText bool, val []byte) error

	Bytes() []byte

	WriteZeros(l int)
}

type encoder struct {
	bs []byte
	i  int
}

func NewEncoder(l int) Encoder {
	return &encoder{bs: make([]byte, l)}
}

func (e *encoder) WriteUint64(val uint64) {
	order.PutUint64(e.bs[e.i:e.i+8], val)
	e.i += 8
}

func (e *encoder) WriteInt64(val int64) {
	order.PutUint64(e.bs[e.i:e.i+8], uint64(val))
	e.i += 8
}

func (e *encoder) WriteUint32(val uint32) {
	order.PutUint32(e.bs[e.i:e.i+4], val)
	e.i += 4
}

func (e *encoder) WriteInt32(val int32) {
	order.PutUint32(e.bs[e.i:e.i+4], uint32(val))
	e.i += 4
}

func (e *encoder) WriteUint16(val uint16) {
	order.PutUint16(e.bs[e.i:e.i+2], val)
	e.i += 2
}

func (e *encoder) WriteInt16(val int16) {
	order.PutUint16(e.bs[e.i:e.i+2], uint16(val))
	e.i += 2
}

func (e *encoder) WriteUint8(val uint8) {
	e.bs[e.i] = val
	e.i += 1
}

func (e *encoder) WriteInt8(val int8) {
	e.bs[e.i] = uint8(val)
	e.i += 1
}

func (e *encoder) WriteBytes(val []byte) {
	l := len(val)
	copy(e.bs[e.i:e.i+l], val)
	e.i += l
}

func (e *encoder) WriteTypeAndInt32Len(isText bool, len int32) {
	if isText {
		e.WriteInt32(len | math.MinInt32)
	}
	e.i += 4
}

func (e *encoder) WriteStringBytesWithInt32Len(isText bool, val []byte) error {
	l := len(val)
	if isText {
		e.WriteInt32(int32(l) | math.MinInt32)
		e.WriteBytes(val)
	} else {
		l /= 2
		e.WriteInt32(int32(l))
		if _, err := hex.Decode(e.bs[e.i:e.i+l], val); err != nil {
			return err
		}
		e.i += l
	}
	return nil
}

func (e *encoder) WriteBytesWithInt32Len(isText bool, val []byte) error {
	l := len(val)
	if isText {
		e.WriteInt32(int32(l) | math.MinInt32)
		e.WriteBytes(val)
	} else {
		e.WriteInt32(int32(l))
		if _, err := hex.Decode(e.bs[e.i:e.i+l], val); err != nil {
			return err
		}
		e.i += l
	}
	return nil
}

func (e *encoder) Bytes() []byte {
	return e.bs
}

func (e *encoder) WriteZeros(l int) {
	e.i += l
}
