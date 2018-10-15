package shabal256

import (
	"hash"
)

// The size of a SHABAL256 checksum in bytes.
const Size = 32

// The blocksize of SHABAL256 in bytes.
const BlockSize = 64

const ivSize = 44

var iv = [ivSize]uint32{
	1392002386,
	3846928793,
	764339180,
	3110359441,
	3758590854,
	3145483465,
	3535126986,
	2966612876,
	349067845,
	581914844,
	4026383467,
	3944855370,
	3042297582,
	1047594390,
	2804573487,
	2466337119,
	3660104186,
	1768937576,
	2629222258,
	184434690,
	2799711765,
	1362674132,
	3189859078,
	3012266128,
	1051244907,
	848932068,
	814894548,
	1439380645,
	3020288049,
	3290644154,
	3010673017,
	3235749205,
	3306956974,
	2737289441,
	1455776103,
	3982574643,
	2293603680,
	1625476794,
	1972063115,
	2213030527,
	3163981864,
	3873442807,
	3129187925,
	2605259872,
}

// digest represents the partial evaluation of a checksum.
type digest struct {
	buf   [64]byte
	state [44]uint32
	ptr   uint32
	w     int64
}

func (d *digest) Reset() {
	copy(d.state[:], iv[:])
	d.w = 1
	d.ptr = 0
}

// New returns a new hash.Hash computing the SHABAL256 checksum.
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (int, error) {
	var off uint32
	nn := len(p)
	len := uint32(nn)

	if d.ptr != 0 {
		rlen := 64 - d.ptr
		if len < rlen {
			copy(d.buf[d.ptr:d.ptr+rlen], p)
			return nn, nil
		}

		off += rlen
		len -= rlen
		d.core1()
	}

	num := len >> 6
	if num > 0 {
		d.core(p, off, num)
		off += num << 6
		len &= 63
	}
	copy(d.buf[:len], p[off:])
	d.ptr = len

	return nn, nil
}

func (d0 *digest) Sum(in []byte) []byte {
	// Make a copy of d0 so that caller can keep writing and summing.
	d := *d0
	hash := d.checkSum()
	return append(in, hash[:]...)
}

func (d *digest) checkSum() [Size]byte {
	var digest [Size]byte
	d.buf[d.ptr] = 0x80
	d.ptr++
	for i := d.ptr; i < BlockSize; i++ {
		d.buf[i] = 0
	}

	d.core1()
	d.w--
	d.core1()
	d.w--
	d.core1()
	d.w--
	d.core1()
	d.w--

	var j int = 36
	var w uint32 = 0

	for i := 0; i < 32; i++ {
		if i&3 == 0 {
			w = d.state[j]
			j++
		}
		digest[i] = byte(w)
		w >>= 8
	}

	return digest
}

func decodeLEInt(data []byte, off uint32) uint32 {
	return uint32(data[off]&0xFF) |
		(uint32(data[off+1]&0xFF) << 8) |
		(uint32(data[off+2]&0xFF) << 16) |
		(uint32(data[off+3]&0xFF) << 24)
}

func (d *digest) core(data []byte, off, num uint32) {
	A0 := d.state[0]
	A1 := d.state[1]
	A2 := d.state[2]
	A3 := d.state[3]
	A4 := d.state[4]
	A5 := d.state[5]
	A6 := d.state[6]
	A7 := d.state[7]
	A8 := d.state[8]
	A9 := d.state[9]
	AA := d.state[10]
	AB := d.state[11]

	B0 := d.state[12]
	B1 := d.state[13]
	B2 := d.state[14]
	B3 := d.state[15]
	B4 := d.state[16]
	B5 := d.state[17]
	B6 := d.state[18]
	B7 := d.state[19]
	B8 := d.state[20]
	B9 := d.state[21]
	BA := d.state[22]
	BB := d.state[23]
	BC := d.state[24]
	BD := d.state[25]
	BE := d.state[26]
	BF := d.state[27]

	C0 := d.state[28]
	C1 := d.state[29]
	C2 := d.state[30]
	C3 := d.state[31]
	C4 := d.state[32]
	C5 := d.state[33]
	C6 := d.state[34]
	C7 := d.state[35]
	C8 := d.state[36]
	C9 := d.state[37]
	CA := d.state[38]
	CB := d.state[39]
	CC := d.state[40]
	CD := d.state[41]
	CE := d.state[42]
	CF := d.state[43]

	for ; num > 0; num-- {
		M0 := decodeLEInt(data, off)

		B0 += M0
		B0 = (B0 << 17) | (B0 >> 15)
		M1 := decodeLEInt(data, off+4)
		B1 += M1
		B1 = (B1 << 17) | (B1 >> 15)
		M2 := decodeLEInt(data, off+8)
		B2 += M2
		B2 = (B2 << 17) | (B2 >> 15)
		M3 := decodeLEInt(data, off+12)
		B3 += M3
		B3 = (B3 << 17) | (B3 >> 15)
		M4 := decodeLEInt(data, off+16)
		B4 += M4
		B4 = (B4 << 17) | (B4 >> 15)
		M5 := decodeLEInt(data, off+20)
		B5 += M5
		B5 = (B5 << 17) | (B5 >> 15)
		M6 := decodeLEInt(data, off+24)
		B6 += M6
		B6 = (B6 << 17) | (B6 >> 15)
		M7 := decodeLEInt(data, off+28)
		B7 += M7
		B7 = (B7 << 17) | (B7 >> 15)
		M8 := decodeLEInt(data, off+32)
		B8 += M8
		B8 = (B8 << 17) | (B8 >> 15)
		M9 := decodeLEInt(data, off+36)
		B9 += M9
		B9 = (B9 << 17) | (B9 >> 15)
		MA := decodeLEInt(data, off+40)
		BA += MA
		BA = (BA << 17) | (BA >> 15)
		MB := decodeLEInt(data, off+44)
		BB += MB
		BB = (BB << 17) | (BB >> 15)
		MC := decodeLEInt(data, off+48)
		BC += MC
		BC = (BC << 17) | (BC >> 15)
		MD := decodeLEInt(data, off+52)
		BD += MD
		BD = (BD << 17) | (BD >> 15)
		ME := decodeLEInt(data, off+56)
		BE += ME
		BE = (BE << 17) | (BE >> 15)
		MF := decodeLEInt(data, off+60)
		BF += MF
		BF = (BF << 17) | (BF >> 15)

		off += 64
		A0 ^= uint32(d.w)
		A1 ^= uint32(d.w >> 32)
		d.w++

		A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ C8) * 3) ^ BD ^ (B9 & ^B6) ^ M0
		B0 = ^((B0 << 1) | (B0 >> 31)) ^ A0
		A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ C7) * 3) ^ BE ^ (BA & ^B7) ^ M1
		B1 = ^((B1 << 1) | (B1 >> 31)) ^ A1
		A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ C6) * 3) ^ BF ^ (BB & ^B8) ^ M2
		B2 = ^((B2 << 1) | (B2 >> 31)) ^ A2
		A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ C5) * 3) ^ B0 ^ (BC & ^B9) ^ M3
		B3 = ^((B3 << 1) | (B3 >> 31)) ^ A3
		A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ C4) * 3) ^ B1 ^ (BD & ^BA) ^ M4
		B4 = ^((B4 << 1) | (B4 >> 31)) ^ A4
		A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ C3) * 3) ^ B2 ^ (BE & ^BB) ^ M5
		B5 = ^((B5 << 1) | (B5 >> 31)) ^ A5
		A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ C2) * 3) ^ B3 ^ (BF & ^BC) ^ M6
		B6 = ^((B6 << 1) | (B6 >> 31)) ^ A6
		A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ C1) * 3) ^ B4 ^ (B0 & ^BD) ^ M7
		B7 = ^((B7 << 1) | (B7 >> 31)) ^ A7
		A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ C0) * 3) ^ B5 ^ (B1 & ^BE) ^ M8

		B8 = ^((B8 << 1) | (B8 >> 31)) ^ A8
		A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ CF) * 3) ^ B6 ^ (B2 & ^BF) ^ M9
		B9 = ^((B9 << 1) | (B9 >> 31)) ^ A9
		AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ CE) * 3) ^ B7 ^ (B3 & ^B0) ^ MA
		BA = ^((BA << 1) | (BA >> 31)) ^ AA
		AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ CD) * 3) ^ B8 ^ (B4 & ^B1) ^ MB
		BB = ^((BB << 1) | (BB >> 31)) ^ AB
		A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ CC) * 3) ^ B9 ^ (B5 & ^B2) ^ MC
		BC = ^((BC << 1) | (BC >> 31)) ^ A0
		A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ CB) * 3) ^ BA ^ (B6 & ^B3) ^ MD
		BD = ^((BD << 1) | (BD >> 31)) ^ A1
		A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ CA) * 3) ^ BB ^ (B7 & ^B4) ^ ME
		BE = ^((BE << 1) | (BE >> 31)) ^ A2
		A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ C9) * 3) ^ BC ^ (B8 & ^B5) ^ MF
		BF = ^((BF << 1) | (BF >> 31)) ^ A3
		A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ C8) * 3) ^ BD ^ (B9 & ^B6) ^ M0
		B0 = ^((B0 << 1) | (B0 >> 31)) ^ A4
		A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ C7) * 3) ^ BE ^ (BA & ^B7) ^ M1
		B1 = ^((B1 << 1) | (B1 >> 31)) ^ A5
		A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ C6) * 3) ^ BF ^ (BB & ^B8) ^ M2
		B2 = ^((B2 << 1) | (B2 >> 31)) ^ A6
		A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ C5) * 3) ^ B0 ^ (BC & ^B9) ^ M3
		B3 = ^((B3 << 1) | (B3 >> 31)) ^ A7
		A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ C4) * 3) ^ B1 ^ (BD & ^BA) ^ M4
		B4 = ^((B4 << 1) | (B4 >> 31)) ^ A8
		A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ C3) * 3) ^ B2 ^ (BE & ^BB) ^ M5
		B5 = ^((B5 << 1) | (B5 >> 31)) ^ A9
		AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ C2) * 3) ^ B3 ^ (BF & ^BC) ^ M6
		B6 = ^((B6 << 1) | (B6 >> 31)) ^ AA
		AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ C1) * 3) ^ B4 ^ (B0 & ^BD) ^ M7
		B7 = ^((B7 << 1) | (B7 >> 31)) ^ AB
		A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ C0) * 3) ^ B5 ^ (B1 & ^BE) ^ M8
		B8 = ^((B8 << 1) | (B8 >> 31)) ^ A0
		A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ CF) * 3) ^ B6 ^ (B2 & ^BF) ^ M9
		B9 = ^((B9 << 1) | (B9 >> 31)) ^ A1
		A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ CE) * 3) ^ B7 ^ (B3 & ^B0) ^ MA
		BA = ^((BA << 1) | (BA >> 31)) ^ A2
		A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ CD) * 3) ^ B8 ^ (B4 & ^B1) ^ MB
		BB = ^((BB << 1) | (BB >> 31)) ^ A3
		A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ CC) * 3) ^ B9 ^ (B5 & ^B2) ^ MC
		BC = ^((BC << 1) | (BC >> 31)) ^ A4
		A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ CB) * 3) ^ BA ^ (B6 & ^B3) ^ MD
		BD = ^((BD << 1) | (BD >> 31)) ^ A5
		A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ CA) * 3) ^ BB ^ (B7 & ^B4) ^ ME
		BE = ^((BE << 1) | (BE >> 31)) ^ A6
		A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ C9) * 3) ^ BC ^ (B8 & ^B5) ^ MF
		BF = ^((BF << 1) | (BF >> 31)) ^ A7
		A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ C8) * 3) ^ BD ^ (B9 & ^B6) ^ M0
		B0 = ^((B0 << 1) | (B0 >> 31)) ^ A8
		A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ C7) * 3) ^ BE ^ (BA & ^B7) ^ M1
		B1 = ^((B1 << 1) | (B1 >> 31)) ^ A9
		AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ C6) * 3) ^ BF ^ (BB & ^B8) ^ M2
		B2 = ^((B2 << 1) | (B2 >> 31)) ^ AA
		AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ C5) * 3) ^ B0 ^ (BC & ^B9) ^ M3
		B3 = ^((B3 << 1) | (B3 >> 31)) ^ AB
		A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ C4) * 3) ^ B1 ^ (BD & ^BA) ^ M4
		B4 = ^((B4 << 1) | (B4 >> 31)) ^ A0
		A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ C3) * 3) ^ B2 ^ (BE & ^BB) ^ M5
		B5 = ^((B5 << 1) | (B5 >> 31)) ^ A1
		A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ C2) * 3) ^ B3 ^ (BF & ^BC) ^ M6
		B6 = ^((B6 << 1) | (B6 >> 31)) ^ A2
		A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ C1) * 3) ^ B4 ^ (B0 & ^BD) ^ M7
		B7 = ^((B7 << 1) | (B7 >> 31)) ^ A3
		A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ C0) * 3) ^ B5 ^ (B1 & ^BE) ^ M8
		B8 = ^((B8 << 1) | (B8 >> 31)) ^ A4
		A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ CF) * 3) ^ B6 ^ (B2 & ^BF) ^ M9
		B9 = ^((B9 << 1) | (B9 >> 31)) ^ A5
		A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ CE) * 3) ^ B7 ^ (B3 & ^B0) ^ MA
		BA = ^((BA << 1) | (BA >> 31)) ^ A6
		A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ CD) * 3) ^ B8 ^ (B4 & ^B1) ^ MB
		BB = ^((BB << 1) | (BB >> 31)) ^ A7
		A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ CC) * 3) ^ B9 ^ (B5 & ^B2) ^ MC
		BC = ^((BC << 1) | (BC >> 31)) ^ A8
		A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ CB) * 3) ^ BA ^ (B6 & ^B3) ^ MD
		BD = ^((BD << 1) | (BD >> 31)) ^ A9
		AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ CA) * 3) ^ BB ^ (B7 & ^B4) ^ ME
		BE = ^((BE << 1) | (BE >> 31)) ^ AA
		AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ C9) * 3) ^ BC ^ (B8 & ^B5) ^ MF
		BF = ^((BF << 1) | (BF >> 31)) ^ AB

		AB += C6 + CA + CE
		AA += C5 + C9 + CD
		A9 += C4 + C8 + CC
		A8 += C3 + C7 + CB
		A7 += C2 + C6 + CA
		A6 += C1 + C5 + C9
		A5 += C0 + C4 + C8
		A4 += CF + C3 + C7
		A3 += CE + C2 + C6
		A2 += CD + C1 + C5
		A1 += CC + C0 + C4
		A0 += CB + CF + C3

		var tmp uint32
		tmp = B0
		B0 = C0 - M0
		C0 = tmp
		tmp = B1
		B1 = C1 - M1
		C1 = tmp
		tmp = B2
		B2 = C2 - M2
		C2 = tmp
		tmp = B3
		B3 = C3 - M3
		C3 = tmp
		tmp = B4
		B4 = C4 - M4
		C4 = tmp
		tmp = B5
		B5 = C5 - M5
		C5 = tmp
		tmp = B6
		B6 = C6 - M6
		C6 = tmp
		tmp = B7
		B7 = C7 - M7
		C7 = tmp
		tmp = B8
		B8 = C8 - M8
		C8 = tmp
		tmp = B9
		B9 = C9 - M9
		C9 = tmp
		tmp = BA
		BA = CA - MA
		CA = tmp
		tmp = BB
		BB = CB - MB
		CB = tmp
		tmp = BC
		BC = CC - MC
		CC = tmp
		tmp = BD
		BD = CD - MD
		CD = tmp
		tmp = BE
		BE = CE - ME
		CE = tmp
		tmp = BF
		BF = CF - MF
		CF = tmp
	}

	d.state[0] = A0
	d.state[1] = A1
	d.state[2] = A2
	d.state[3] = A3
	d.state[4] = A4
	d.state[5] = A5
	d.state[6] = A6
	d.state[7] = A7
	d.state[8] = A8
	d.state[9] = A9
	d.state[10] = AA
	d.state[11] = AB

	d.state[12] = B0
	d.state[13] = B1
	d.state[14] = B2
	d.state[15] = B3
	d.state[16] = B4
	d.state[17] = B5
	d.state[18] = B6
	d.state[19] = B7
	d.state[20] = B8
	d.state[21] = B9
	d.state[22] = BA
	d.state[23] = BB
	d.state[24] = BC
	d.state[25] = BD
	d.state[26] = BE
	d.state[27] = BF

	d.state[28] = C0
	d.state[29] = C1
	d.state[30] = C2
	d.state[31] = C3
	d.state[32] = C4
	d.state[33] = C5
	d.state[34] = C6
	d.state[35] = C7
	d.state[36] = C8
	d.state[37] = C9
	d.state[38] = CA
	d.state[39] = CB
	d.state[40] = CC
	d.state[41] = CD
	d.state[42] = CE
	d.state[43] = CF
}

func (d *digest) core1() {
	A0 := d.state[0]
	A1 := d.state[1]
	A2 := d.state[2]
	A3 := d.state[3]
	A4 := d.state[4]
	A5 := d.state[5]
	A6 := d.state[6]
	A7 := d.state[7]
	A8 := d.state[8]
	A9 := d.state[9]
	AA := d.state[10]
	AB := d.state[11]

	B0 := d.state[12]
	B1 := d.state[13]
	B2 := d.state[14]
	B3 := d.state[15]
	B4 := d.state[16]
	B5 := d.state[17]
	B6 := d.state[18]
	B7 := d.state[19]
	B8 := d.state[20]
	B9 := d.state[21]
	BA := d.state[22]
	BB := d.state[23]
	BC := d.state[24]
	BD := d.state[25]
	BE := d.state[26]
	BF := d.state[27]

	C0 := d.state[28]
	C1 := d.state[29]
	C2 := d.state[30]
	C3 := d.state[31]
	C4 := d.state[32]
	C5 := d.state[33]
	C6 := d.state[34]
	C7 := d.state[35]
	C8 := d.state[36]
	C9 := d.state[37]
	CA := d.state[38]
	CB := d.state[39]
	CC := d.state[40]
	CD := d.state[41]
	CE := d.state[42]
	CF := d.state[43]

	M0 := decodeLEInt(d.buf[:], 0)
	B0 += M0
	B0 = (B0 << 17) | (B0 >> 15)
	M1 := decodeLEInt(d.buf[:], 4)
	B1 += M1
	B1 = (B1 << 17) | (B1 >> 15)
	M2 := decodeLEInt(d.buf[:], 8)
	B2 += M2
	B2 = (B2 << 17) | (B2 >> 15)
	M3 := decodeLEInt(d.buf[:], 12)
	B3 += M3
	B3 = (B3 << 17) | (B3 >> 15)
	M4 := decodeLEInt(d.buf[:], 16)
	B4 += M4
	B4 = (B4 << 17) | (B4 >> 15)
	M5 := decodeLEInt(d.buf[:], 20)
	B5 += M5
	B5 = (B5 << 17) | (B5 >> 15)
	M6 := decodeLEInt(d.buf[:], 24)
	B6 += M6
	B6 = (B6 << 17) | (B6 >> 15)
	M7 := decodeLEInt(d.buf[:], 28)
	B7 += M7
	B7 = (B7 << 17) | (B7 >> 15)
	M8 := decodeLEInt(d.buf[:], 32)
	B8 += M8
	B8 = (B8 << 17) | (B8 >> 15)
	M9 := decodeLEInt(d.buf[:], 36)
	B9 += M9
	B9 = (B9 << 17) | (B9 >> 15)
	MA := decodeLEInt(d.buf[:], 40)
	BA += MA
	BA = (BA << 17) | (BA >> 15)
	MB := decodeLEInt(d.buf[:], 44)
	BB += MB
	BB = (BB << 17) | (BB >> 15)
	MC := decodeLEInt(d.buf[:], 48)
	BC += MC
	BC = (BC << 17) | (BC >> 15)
	MD := decodeLEInt(d.buf[:], 52)
	BD += MD
	BD = (BD << 17) | (BD >> 15)
	ME := decodeLEInt(d.buf[:], 56)
	BE += ME
	BE = (BE << 17) | (BE >> 15)
	MF := decodeLEInt(d.buf[:], 60)
	BF += MF
	BF = (BF << 17) | (BF >> 15)

	A0 ^= uint32(d.w)
	A1 ^= uint32(d.w >> 32)
	d.w++

	A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ C8) * 3) ^ BD ^ (B9 & ^B6) ^ M0
	B0 = ^((B0 << 1) | (B0 >> 31)) ^ A0
	A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ C7) * 3) ^ BE ^ (BA & ^B7) ^ M1
	B1 = ^((B1 << 1) | (B1 >> 31)) ^ A1
	A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ C6) * 3) ^ BF ^ (BB & ^B8) ^ M2
	B2 = ^((B2 << 1) | (B2 >> 31)) ^ A2
	A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ C5) * 3) ^ B0 ^ (BC & ^B9) ^ M3
	B3 = ^((B3 << 1) | (B3 >> 31)) ^ A3
	A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ C4) * 3) ^ B1 ^ (BD & ^BA) ^ M4
	B4 = ^((B4 << 1) | (B4 >> 31)) ^ A4
	A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ C3) * 3) ^ B2 ^ (BE & ^BB) ^ M5
	B5 = ^((B5 << 1) | (B5 >> 31)) ^ A5
	A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ C2) * 3) ^ B3 ^ (BF & ^BC) ^ M6
	B6 = ^((B6 << 1) | (B6 >> 31)) ^ A6
	A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ C1) * 3) ^ B4 ^ (B0 & ^BD) ^ M7
	B7 = ^((B7 << 1) | (B7 >> 31)) ^ A7
	A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ C0) * 3) ^ B5 ^ (B1 & ^BE) ^ M8

	B8 = ^((B8 << 1) | (B8 >> 31)) ^ A8
	A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ CF) * 3) ^ B6 ^ (B2 & ^BF) ^ M9
	B9 = ^((B9 << 1) | (B9 >> 31)) ^ A9
	AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ CE) * 3) ^ B7 ^ (B3 & ^B0) ^ MA
	BA = ^((BA << 1) | (BA >> 31)) ^ AA
	AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ CD) * 3) ^ B8 ^ (B4 & ^B1) ^ MB
	BB = ^((BB << 1) | (BB >> 31)) ^ AB
	A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ CC) * 3) ^ B9 ^ (B5 & ^B2) ^ MC
	BC = ^((BC << 1) | (BC >> 31)) ^ A0
	A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ CB) * 3) ^ BA ^ (B6 & ^B3) ^ MD
	BD = ^((BD << 1) | (BD >> 31)) ^ A1
	A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ CA) * 3) ^ BB ^ (B7 & ^B4) ^ ME
	BE = ^((BE << 1) | (BE >> 31)) ^ A2
	A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ C9) * 3) ^ BC ^ (B8 & ^B5) ^ MF
	BF = ^((BF << 1) | (BF >> 31)) ^ A3
	A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ C8) * 3) ^ BD ^ (B9 & ^B6) ^ M0
	B0 = ^((B0 << 1) | (B0 >> 31)) ^ A4
	A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ C7) * 3) ^ BE ^ (BA & ^B7) ^ M1
	B1 = ^((B1 << 1) | (B1 >> 31)) ^ A5
	A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ C6) * 3) ^ BF ^ (BB & ^B8) ^ M2
	B2 = ^((B2 << 1) | (B2 >> 31)) ^ A6
	A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ C5) * 3) ^ B0 ^ (BC & ^B9) ^ M3
	B3 = ^((B3 << 1) | (B3 >> 31)) ^ A7
	A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ C4) * 3) ^ B1 ^ (BD & ^BA) ^ M4
	B4 = ^((B4 << 1) | (B4 >> 31)) ^ A8
	A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ C3) * 3) ^ B2 ^ (BE & ^BB) ^ M5
	B5 = ^((B5 << 1) | (B5 >> 31)) ^ A9
	AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ C2) * 3) ^ B3 ^ (BF & ^BC) ^ M6
	B6 = ^((B6 << 1) | (B6 >> 31)) ^ AA
	AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ C1) * 3) ^ B4 ^ (B0 & ^BD) ^ M7
	B7 = ^((B7 << 1) | (B7 >> 31)) ^ AB
	A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ C0) * 3) ^ B5 ^ (B1 & ^BE) ^ M8
	B8 = ^((B8 << 1) | (B8 >> 31)) ^ A0
	A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ CF) * 3) ^ B6 ^ (B2 & ^BF) ^ M9
	B9 = ^((B9 << 1) | (B9 >> 31)) ^ A1
	A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ CE) * 3) ^ B7 ^ (B3 & ^B0) ^ MA
	BA = ^((BA << 1) | (BA >> 31)) ^ A2
	A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ CD) * 3) ^ B8 ^ (B4 & ^B1) ^ MB
	BB = ^((BB << 1) | (BB >> 31)) ^ A3
	A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ CC) * 3) ^ B9 ^ (B5 & ^B2) ^ MC
	BC = ^((BC << 1) | (BC >> 31)) ^ A4
	A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ CB) * 3) ^ BA ^ (B6 & ^B3) ^ MD
	BD = ^((BD << 1) | (BD >> 31)) ^ A5
	A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ CA) * 3) ^ BB ^ (B7 & ^B4) ^ ME
	BE = ^((BE << 1) | (BE >> 31)) ^ A6
	A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ C9) * 3) ^ BC ^ (B8 & ^B5) ^ MF
	BF = ^((BF << 1) | (BF >> 31)) ^ A7
	A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ C8) * 3) ^ BD ^ (B9 & ^B6) ^ M0
	B0 = ^((B0 << 1) | (B0 >> 31)) ^ A8
	A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ C7) * 3) ^ BE ^ (BA & ^B7) ^ M1
	B1 = ^((B1 << 1) | (B1 >> 31)) ^ A9
	AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ C6) * 3) ^ BF ^ (BB & ^B8) ^ M2
	B2 = ^((B2 << 1) | (B2 >> 31)) ^ AA
	AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ C5) * 3) ^ B0 ^ (BC & ^B9) ^ M3
	B3 = ^((B3 << 1) | (B3 >> 31)) ^ AB
	A0 = ((A0 ^ (((AB << 15) | (AB >> 17)) * 5) ^ C4) * 3) ^ B1 ^ (BD & ^BA) ^ M4
	B4 = ^((B4 << 1) | (B4 >> 31)) ^ A0
	A1 = ((A1 ^ (((A0 << 15) | (A0 >> 17)) * 5) ^ C3) * 3) ^ B2 ^ (BE & ^BB) ^ M5
	B5 = ^((B5 << 1) | (B5 >> 31)) ^ A1
	A2 = ((A2 ^ (((A1 << 15) | (A1 >> 17)) * 5) ^ C2) * 3) ^ B3 ^ (BF & ^BC) ^ M6
	B6 = ^((B6 << 1) | (B6 >> 31)) ^ A2
	A3 = ((A3 ^ (((A2 << 15) | (A2 >> 17)) * 5) ^ C1) * 3) ^ B4 ^ (B0 & ^BD) ^ M7
	B7 = ^((B7 << 1) | (B7 >> 31)) ^ A3
	A4 = ((A4 ^ (((A3 << 15) | (A3 >> 17)) * 5) ^ C0) * 3) ^ B5 ^ (B1 & ^BE) ^ M8
	B8 = ^((B8 << 1) | (B8 >> 31)) ^ A4
	A5 = ((A5 ^ (((A4 << 15) | (A4 >> 17)) * 5) ^ CF) * 3) ^ B6 ^ (B2 & ^BF) ^ M9
	B9 = ^((B9 << 1) | (B9 >> 31)) ^ A5
	A6 = ((A6 ^ (((A5 << 15) | (A5 >> 17)) * 5) ^ CE) * 3) ^ B7 ^ (B3 & ^B0) ^ MA
	BA = ^((BA << 1) | (BA >> 31)) ^ A6
	A7 = ((A7 ^ (((A6 << 15) | (A6 >> 17)) * 5) ^ CD) * 3) ^ B8 ^ (B4 & ^B1) ^ MB
	BB = ^((BB << 1) | (BB >> 31)) ^ A7
	A8 = ((A8 ^ (((A7 << 15) | (A7 >> 17)) * 5) ^ CC) * 3) ^ B9 ^ (B5 & ^B2) ^ MC
	BC = ^((BC << 1) | (BC >> 31)) ^ A8
	A9 = ((A9 ^ (((A8 << 15) | (A8 >> 17)) * 5) ^ CB) * 3) ^ BA ^ (B6 & ^B3) ^ MD
	BD = ^((BD << 1) | (BD >> 31)) ^ A9
	AA = ((AA ^ (((A9 << 15) | (A9 >> 17)) * 5) ^ CA) * 3) ^ BB ^ (B7 & ^B4) ^ ME
	BE = ^((BE << 1) | (BE >> 31)) ^ AA
	AB = ((AB ^ (((AA << 15) | (AA >> 17)) * 5) ^ C9) * 3) ^ BC ^ (B8 & ^B5) ^ MF
	BF = ^((BF << 1) | (BF >> 31)) ^ AB

	AB += C6 + CA + CE
	AA += C5 + C9 + CD
	A9 += C4 + C8 + CC
	A8 += C3 + C7 + CB
	A7 += C2 + C6 + CA
	A6 += C1 + C5 + C9
	A5 += C0 + C4 + C8
	A4 += CF + C3 + C7
	A3 += CE + C2 + C6
	A2 += CD + C1 + C5
	A1 += CC + C0 + C4
	A0 += CB + CF + C3

	var tmp uint32
	tmp = B0
	B0 = C0 - M0
	C0 = tmp
	tmp = B1
	B1 = C1 - M1
	C1 = tmp
	tmp = B2
	B2 = C2 - M2
	C2 = tmp
	tmp = B3
	B3 = C3 - M3
	C3 = tmp
	tmp = B4
	B4 = C4 - M4
	C4 = tmp
	tmp = B5
	B5 = C5 - M5
	C5 = tmp
	tmp = B6
	B6 = C6 - M6
	C6 = tmp
	tmp = B7
	B7 = C7 - M7
	C7 = tmp
	tmp = B8
	B8 = C8 - M8
	C8 = tmp
	tmp = B9
	B9 = C9 - M9
	C9 = tmp
	tmp = BA
	BA = CA - MA
	CA = tmp
	tmp = BB
	BB = CB - MB
	CB = tmp
	tmp = BC
	BC = CC - MC
	CC = tmp
	tmp = BD
	BD = CD - MD
	CD = tmp
	tmp = BE
	BE = CE - ME
	CE = tmp
	tmp = BF
	BF = CF - MF
	CF = tmp

	d.state[0] = A0
	d.state[1] = A1
	d.state[2] = A2
	d.state[3] = A3
	d.state[4] = A4
	d.state[5] = A5
	d.state[6] = A6
	d.state[7] = A7
	d.state[8] = A8
	d.state[9] = A9
	d.state[10] = AA
	d.state[11] = AB

	d.state[12] = B0
	d.state[13] = B1
	d.state[14] = B2
	d.state[15] = B3
	d.state[16] = B4
	d.state[17] = B5
	d.state[18] = B6
	d.state[19] = B7
	d.state[20] = B8
	d.state[21] = B9
	d.state[22] = BA
	d.state[23] = BB
	d.state[24] = BC
	d.state[25] = BD
	d.state[26] = BE
	d.state[27] = BF

	d.state[28] = C0
	d.state[29] = C1
	d.state[30] = C2
	d.state[31] = C3
	d.state[32] = C4
	d.state[33] = C5
	d.state[34] = C6
	d.state[35] = C7
	d.state[36] = C8
	d.state[37] = C9
	d.state[38] = CA
	d.state[39] = CB
	d.state[40] = CC
	d.state[41] = CD
	d.state[42] = CE
	d.state[43] = CF
}

// Sum256 returns the SHABAL256 checksum of the data.
func Sum256(data []byte) [Size]byte {
	var d digest
	d.Reset()
	d.Write(data)
	return d.checkSum()
}
