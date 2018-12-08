package encoding

type Decoder interface {
	ReadUint64() uint64
	ReadInt64() int64
	ReadUint32() uint32
	ReadInt32() int32
	ReadUint16() uint16
	ReadInt16() int16
	ReadUint8() uint8
	ReadInt8() int8

	ReadBytes(l int) []byte

	Step(i int)
	Position() int
	Reset(i int)
}

type decoder struct {
	bs []byte
	i  int
}

func NewDecoder(bs []byte) Decoder {
	return &decoder{bs: bs}
}

func (d *decoder) ReadUint64() uint64 {
	val := order.Uint64(d.bs[d.i : d.i+8])
	d.i += 8
	return val
}

func (d *decoder) ReadInt64() int64 {
	val := order.Uint64(d.bs[d.i : d.i+8])
	d.i += 8
	return int64(val)
}

func (d *decoder) ReadUint32() uint32 {
	val := order.Uint32(d.bs[d.i : d.i+4])
	d.i += 4
	return val
}

func (d *decoder) ReadInt32() int32 {
	val := order.Uint32(d.bs[d.i : d.i+4])
	d.i += 4
	return int32(val)
}

func (d *decoder) ReadUint16() uint16 {
	val := order.Uint16(d.bs[d.i : d.i+2])
	d.i += 2
	return val
}

func (d *decoder) ReadInt16() int16 {
	val := order.Uint16(d.bs[d.i : d.i+2])
	d.i += 2
	return int16(val)
}

func (d *decoder) ReadUint8() uint8 {
	val := d.bs[d.i]
	d.i += 1
	return val
}

func (d *decoder) ReadInt8() int8 {
	val := d.bs[d.i]
	d.i += 1
	return int8(val)
}

func (d *decoder) ReadBytes(l int) []byte {
	bs := d.bs[d.i : d.i+l]
	d.i += l
	return bs
}

func (d *decoder) Step(i int) {
	d.i += i
}

func (d *decoder) Position() int {
	return d.i
}

func (d *decoder) Reset(i int) {
	d.i = i
}
