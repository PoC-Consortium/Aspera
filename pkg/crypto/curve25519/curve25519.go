package curve25519

var order = []byte{
	237, 211, 245, 92, 26, 99, 18, 88, 214, 156, 247, 162, 222, 249, 222, 20, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 16,
}

func clamp(k *[32]byte) {
	k[31] &= 0x7F
	k[31] |= 0x40
	k[0] &= 0xF8
}

func mulaSmall(p, q []byte, m int, x []byte, n, z int) int {
	v := 0
	for i := 0; i < n; i++ {
		v += int(q[i+m]) + z*int(x[i])
		p[i+m] = byte(v)
		v >>= 8
	}
	return v
}

func mula32(p, x, y []byte, t, z int) int {
	n := 31
	var w, i int
	for ; i < t; i++ {
		zy := z * int(y[i])
		w += mulaSmall(p, p, i, x, n, zy) + int(p[i+n]) + zy*int(x[n])
		p[i+n] = byte(w)
		w >>= 8
	}
	p[i+n] = byte(w + int(p[i+n]))
	return w >> 8
}

func divmod(q, r []byte, n int, d []byte, t int) {
	rn := 0
	dt := int(d[t-1]) << 8
	if t > 1 {
		dt |= int(d[t-2])
	}
	for n--; n >= t-1; n-- {
		z := (rn << 16) | (int(r[n]) << 8)
		if n > 0 {
			z |= int(r[n-1])
		}
		z /= dt
		rn += mulaSmall(r, r, n-t+1, d, t, -z)
		q[n-t+1] = byte(z + rn)
		mulaSmall(r, r, n-t+1, d, t, -rn)
		rn = int(r[n])
		// if n == 15 {
		// 	log.Fatal(r)
		// }
		r[n] = 0
	}
	r[t-1] = byte(rn)
}

func unpack(x *[10]int64, m []byte) {
	x[0] = int64(m[0]) | (int64(m[1]) << 8) | (int64(m[2]) << 16) | ((int64(m[3]) & 3) << 24)
	x[1] = ((int64(m[3]) & ^3) >> 2) | (int64(m[4]) << 6) | (int64(m[5]) << 14) | (((int64(m[6]) & 0xFF) & 7) << 22)
	x[2] = ((int64(m[6]) & ^7) >> 3) | (int64(m[7]) << 5) | (int64(m[8]) << 13) | ((int64(m[9]) & 31) << 21)
	x[3] = ((int64(m[9]) & ^31) >> 5) | (int64(m[10]) << 3) | ((int64(m[11]) & 0xFF) << 11) | ((int64(m[12]) & 63) << 19)
	x[4] = (((int64(m[12]) & 0xFF) & ^63) >> 6) | (int64(m[13]) << 2) | (int64(m[14]) << 10) | (int64(m[15]) << 18)
	x[5] = int64(m[16]) | (int64(m[17]) << 8) | (int64(m[18]) << 16) | ((int64(m[19]) & 1) << 24)
	x[6] = (((int64(m[19]) & 0xFF) & ^1) >> 1) | (int64(m[20]) << 7) | (int64(m[21]) << 15) | ((int64(m[22]) & 7) << 23)
	x[7] = ((int64(m[22]) & ^7) >> 3) | (int64(m[23]) << 5) | (int64(m[24]) << 13) | ((int64(m[25]) & 15) << 21)
	x[8] = ((int64(m[25]) & ^15) >> 4) | (int64(m[26]) << 4) | (int64(m[27]) << 12) | ((int64(m[28]) & 63) << 20)
	x[9] = ((int64(m[28]) & ^63) >> 6) | (int64(m[29]) << 2) | (int64(m[30]) << 10) | (int64(m[31]) << 18)
}

func sign(v, h, x, s []byte) bool {
	var w int
	h1 := make([]byte, 32)
	x1 := make([]byte, 32)
	tmp3 := make([]byte, 32)
	tmp1 := make([]byte, 64)
	tmp2 := make([]byte, 64)

	copy(h1, h)
	copy(x1, x)

	divmod(tmp3, h1, 32, order, 32)
	divmod(tmp3, x1, 32, order, 32)

	mulaSmall(v, x1, 0, h1, 32, -1)
	mulaSmall(v, v, 0, order, 32, 1)

	mula32(tmp1, v, s, 32, 1)
	divmod(tmp2, tmp1, 64, order, 32)

	for i := 0; i < 32; i++ {
		v[i] = tmp1[i]
		w |= int(tmp1[i])
	}
	return w != 0
}
