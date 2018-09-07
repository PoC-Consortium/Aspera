package curve25519

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClamp(t *testing.T) {
	bs := make([]byte, 32)
	bs[31] = 0x64
	bs[31] = 0x04
	bs[0] = 0xE5
	clamp(bs)
	assert.Equal(t, byte(0x44), bs[31])
	assert.Equal(t, byte(0xe0), bs[0])
}

func initByteSlice(len, offset int) []byte {
	bs := make([]byte, len)
	for i := range bs {
		bs[i] = byte(offset + i)
	}
	return bs
}

func TestMulaSmall(t *testing.T) {
	p := initByteSlice(32, 0)
	q := initByteSlice(32, 32)
	x := initByteSlice(32, 64)
	v := mulaSmall(p, q, 4, x, 28, 37)
	pExp := []byte{
		0, 1, 2, 3, 100, 147, 185, 223,
		5, 44, 82, 120, 158, 196, 234, 16,
		55, 93, 131, 169, 207, 245, 27, 66,
		104, 142, 180, 218, 0, 39, 77, 115,
	}
	assert.Equal(t, 13, v)
	assert.Equal(t, pExp, p)
}

func TestMula32(t *testing.T) {
	p := initByteSlice(48, 0)
	x := initByteSlice(32, 32)
	y := initByteSlice(16, 64)
	mula32(p, x, y, 16, 237)
	pExp := []byte{
		0, 73, 191, 73, 214, 82, 173, 211,
		179, 59, 89, 250, 12, 127, 62, 57,
		93, 166, 230, 38, 103, 167, 231, 39,
		104, 168, 232, 40, 105, 169, 233, 41,
		106, 90, 88, 132, 240, 174, 209, 106,
		140, 72, 177, 216, 208, 171, 123, 82,
	}
	assert.Equal(t, pExp, p)
}

func TestDivmod(t *testing.T) {
	q := initByteSlice(17, 0)
	r := initByteSlice(32, 32)
	d := initByteSlice(16, 64)
	divmod(q, r, 32, d, 16)
	rExp := []byte{
		224, 15, 24, 169, 151, 34, 235, 230,
		238, 66, 181, 168, 235, 89, 137, 23,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	qExp := []byte{
		65, 198, 83, 44, 242, 203, 255, 194,
		131, 104, 122, 114, 39, 58, 38, 204,
		0,
	}
	assert.Equal(t, rExp, r)
	assert.Equal(t, qExp, q)
}

func TestUnpack(t *testing.T) {
	m := make([]byte, 32)
	for i := range m {
		m[i] = byte(i)
	}
	var x long10
	unpack(&x, m)
	assert.Equal(t, int64(50462976), x[0])
	assert.Equal(t, int64(25248000), x[1])
	assert.Equal(t, int64(18940128), x[2])
	assert.Equal(t, int64(6314064), x[3])
	assert.Equal(t, int64(3946548), x[4])
	assert.Equal(t, int64(17961232), x[5])
	assert.Equal(t, int64(51022345), x[6])
	assert.Equal(t, int64(19071714), x[7])
	assert.Equal(t, int64(29471137), x[8])
	assert.Equal(t, int64(8157300), x[9])
}

func TestPack(t *testing.T) {
	m := make([]byte, 32)
	x := long10{
		12973021,
		943701274,
		971409174,
		750923750,
		89384184,
		164918401,
		474091741,
		913461231,
		174917410,
		143091724,
	}
	pack(&x, m)
	mExp := []byte{
		240, 243, 197, 104, 244, 254, 144, 57,
		52, 143, 254, 69, 152, 67, 249, 84,
		130, 116, 212, 195, 33, 132, 176, 127,
		146, 211, 115, 208, 166, 3, 218, 161,
	}
	assert.Equal(t, mExp, m)
}

func TestSign(t *testing.T) {
	v := initByteSlice(32, 0)
	h := initByteSlice(32, 32)
	x := initByteSlice(32, 64)
	s := initByteSlice(32, 96)
	assert.True(t, Sign(v, h, x, s))
	vExp := []byte{
		81, 116, 125, 136, 11, 176, 40, 13,
		125, 197, 85, 64, 246, 64, 162, 63,
		190, 226, 217, 208, 233, 9, 199, 57,
		103, 247, 248, 228, 250, 158, 78, 4,
	}
	assert.Equal(t, vExp, v)
}

func TestNumsize(t *testing.T) {
	x := make([]byte, 32)
	x[13] = 1
	assert.Equal(t, 0, numsize(x, 12))
	assert.Equal(t, 0, numsize(x, 13))
	assert.Equal(t, 14, numsize(x, 14))

	x = make([]byte, 32)
	x[0] = 32
	assert.Equal(t, 1, numsize(x, 32))
}

func TestEgcd32(t *testing.T) {
	x := initByteSlice(64, 0)
	y := initByteSlice(64, 32)
	a := initByteSlice(32, 64)
	b := initByteSlice(32, 96)
	gcdExp := []byte{
		244, 3, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		249, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
	}
	assert.Equal(t, gcdExp, egcd32(x, y, a, b))
}
