package at

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"unsafe"

	"github.com/ac0v/aspera/pkg/common/math"
)

const (
	codePageBytes = 512
	dataPageBytes = 512

	callStackePageBytes = 256
	userStackPageBytes  = 256

	defaultBalance = 100

	maxToMultiply = 0x1fffffff

	opCodeNop           = 0x7f
	opCodeSetVal        = 0x01
	opCodeSetDat        = 0x02
	opCodeClrDat        = 0x03
	opCodeIncDat        = 0x04
	opCodeDecDat        = 0x05
	opCodeAddDat        = 0x06
	opCodeSubDat        = 0x07
	opCodeMulDat        = 0x08
	opCodeDivDat        = 0x09
	opCodeBorDat        = 0x0a
	opCodeAndDat        = 0x0b
	opCodeXorDat        = 0x0c
	opCodeNotDat        = 0x0d
	opCodeSetInd        = 0x0e
	opCodeSetIdx        = 0x0f
	opCodePshDat        = 0x10
	opCodePopDat        = 0x11
	opCodeJmpSub        = 0x12
	opCodeRetSub        = 0x13
	opCodeIndDat        = 0x14
	opCodeIdxDat        = 0x15
	opCodeModDat        = 0x16
	opCodeShlDat        = 0x17
	opCodeShrDat        = 0x18
	opCodeJmpAdr        = 0x1a
	opCodeBzrDat        = 0x1b
	opCodeBnzDat        = 0x1e
	opCodeBgtDat        = 0x1f
	opCodeBltDat        = 0x20
	opCodeBgeDat        = 0x21
	opCodeBleDat        = 0x22
	opCodeBeqDat        = 0x23
	opCodeBneDat        = 0x24
	opCodeSlpDat        = 0x25
	opCodeFizDat        = 0x26
	opCodeStzDat        = 0x27
	opCodeFinImd        = 0x28
	opCodeStpImd        = 0x29
	opCodeSlpImd        = 0x2a
	opCodeErrAdr        = 0x2b
	opCodeSetPcs        = 0x30
	opCodeExtFun        = 0x32
	opCodeExtFunDat     = 0x33
	opCodeExtFunDat2    = 0x34
	opCodeExtFunRet     = 0x35
	opCodeExtFunRetDat  = 0x36
	opCodeExtFunRetDat2 = 0x37

	funGetSize      = 0x0003
	funGetA1        = 0x0100
	funGetA2        = 0x0101
	funGetA3        = 0x0102
	funGetA4        = 0x0103
	funGetB1        = 0x0104
	funGetB2        = 0x0105
	funGetB3        = 0x0106
	funGetB4        = 0x0107
	funClearA       = 0x0120
	funClearB       = 0x0121
	funClearAB      = 0x0122
	funSetAToB      = 0x0123
	funSetBToA      = 0x0124
	funCheckAIsZero = 0x0125
	funCheckBIsZero = 0x0126
	funCheckAIsB    = 0x0127
	funSwapAB       = 0x0128
	funAOrWithB     = 0x0129
	funBOrWithA     = 0x012A
	funAAndWithB    = 0x012b
	funBAndWithA    = 0x012c
	funAXorWithB    = 0x012d
	funBXorWithA    = 0x012e
	funAddAToB      = 0x0140
	funAddBToA      = 0x0141
	funSubAFromB    = 0x0142
	funSubBFromA    = 0x0143
	funMulAByB      = 0x0144
	funMulBByA      = 0x0145
	funDivAByB      = 0x0146
	funDivBByA      = 0x0147

	ok                     = 0
	errCodeOverflow        = -1
	errCodeInvalidCode     = -2
	errCodeUnexpectedError = -3
)

var ErrOverflow = errors.New("overflow")
var ErrInvalidCode = errors.New("invalid code")
var ErrUnexpectedError = errors.New("unexpected error")

var codePages = 1
var dataPages = 1

var callStackPages = 1
var userStackPages = 1

var balance int64 = defaultBalance

var firstCall = true

var incrementFunc int32

type stateMachine struct {
	code []byte
	data []byte

	csize int32
	dsize int32

	stopped  bool
	finished bool

	pc int32

	pce int32
	pcs int32

	opc int32

	cs int32
	us int32

	as [32]byte
	bs [32]byte

	steps int32

	sleepUntil int32

	val int64
}

func newStateMachine(code []byte) *stateMachine {
	codeBuf := make([]byte, codePages*codePageBytes)
	copy(codeBuf, code)

	dsize := dataPages * dataPageBytes
	data := make([]byte, dsize+
		callStackPages*callStackePageBytes+
		userStackPages*userStackPageBytes)

	return &stateMachine{
		code:  codeBuf,
		data:  data,
		csize: int32(len(codeBuf)),
		dsize: int32(dsize),
	}
}

func (s *stateMachine) getData() []byte {
	return s.data[:s.dsize]
}

type functionData struct {
	loop   bool
	offset int32

	data []int64
}

var funData = map[int32]*functionData{}

func getFunctionData(funNum int32) int64 {
	if funNum == incrementFunc {
		if firstCall {
			firstCall = false
		} else {
			for _, f := range funData {
				f.offset++
				if f.offset >= int32(len(f.data)) {
					if f.loop {
						f.offset = 0
					} else {
						f.offset--
					}
				}
			}
		}
	}

	return funData[funNum].data[funData[funNum].offset]
}

func (s *stateMachine) getA(i int) int64 {
	i *= 8
	return int64(binary.LittleEndian.Uint64(s.as[i : i+8]))
}

func (s *stateMachine) getB(i int) int64 {
	i *= 8
	return int64(binary.LittleEndian.Uint64(s.bs[i : i+8]))
}

func (s *stateMachine) setA(i int, val int64) {
	i *= 8
	binary.LittleEndian.PutUint64(s.as[i:i+8], uint64(val))
}

func (s *stateMachine) setB(i int, val int64) {
	i *= 8
	binary.LittleEndian.PutUint64(s.bs[i:i+8], uint64(val))
}

func bytesToSignedBigInt(bs []byte) *big.Int {
	var i big.Int
	i.SetBytes(bs)
	return math.S256(&i)
}

func signedBigIntToBytes(i *big.Int) []byte {
	bs := math.U256(i).Bytes()
	swapEndianInline(bs)
	return bs
}

func (s *stateMachine) getBigAB() (*big.Int, *big.Int) {
	a := bytesToSignedBigInt(swapEndian(s.as[:]))
	b := bytesToSignedBigInt(swapEndian(s.bs[:]))
	return a, b
}

func swapEndian(bs []byte) []byte {
	swapped := make([]byte, len(bs))
	for i := 0; i < len(bs)/2; i++ {
		swapped[i], swapped[len(bs)-i-1] = bs[len(bs)-i-1], bs[i]
	}
	return swapped
}

func swapEndianInline(bs []byte) {
	for i := 0; i < len(bs)/2; i++ {
		bs[i], bs[len(bs)-i-1] = bs[len(bs)-i-1], bs[i]
	}
}

func bigIntToPaddedBuffer(i *big.Int) []byte {
	bs := signedBigIntToBytes(i)
	if i.Sign() == -1 {
		bs[len(bs)-1] |= 0x80
	}

	if len(bs) > 32 {
		bs = bs[:32]
	}
	bytesToPad := 8*4 - len(bs)
	if bytesToPad > 0 {
		var padding byte
		if len(bs) > 0 {
			padding = (bs[0] & 0x80) >> 7
		}
		paddingBytes := make([]byte, bytesToPad)
		for i := 0; i < bytesToPad; i++ {
			paddingBytes[i] = padding
		}
		bs = append(bs, paddingBytes...)
	}
	return bs
}

func (s *stateMachine) setAs(bs []byte) {
	copy(s.as[:], bs)
}

func (s *stateMachine) setBs(bs []byte) {
	copy(s.bs[:], bs)
}

func (s *stateMachine) fun(funNum int32) int64 {
	var rc int64

	switch funNum {
	case 1:
		rc = s.val
	case 2:
		if s.val == 9 {
			rc = 0
			s.val = 0
		} else {
			s.val++
			rc = s.val
		}
	case funGetSize:
		rc = 10
	case 4:
		rc = int64(funNum)
	case 25:
		rc = balance
		if _, exists := funData[funNum]; exists {
			for _, f := range funData {
				f.offset = 0
			}
		}
	case 32:
		rc = balance
		if _, exists := funData[funNum]; exists {
			for _, f := range funData {
				f.offset = 0
			}
		}
	case funGetA1:
		rc = s.getA(0)
	case funGetA2:
		rc = s.getA(1)
	case funGetA3:
		rc = s.getA(2)
	case funGetA4:
		rc = s.getA(3)
	case funGetB1:
		rc = s.getB(0)
	case funGetB2:
		rc = s.getB(1)
	case funGetB3:
		rc = s.getB(2)
	case funGetB4:
		rc = s.getB(3)
	case funClearA:
		for i := range s.as {
			s.as[i] = 0
		}
	case funClearB:
		for i := range s.bs {
			s.bs[i] = 0
		}
	case funClearAB:
		for i := range s.as {
			s.as[i] = 0
			s.bs[i] = 0
		}
	case funSetAToB:
		for i, b := range s.bs {
			s.as[i] = b
		}
	case funSetBToA:
		for i, a := range s.as {
			s.bs[i] = a
		}
	case funCheckAIsZero:
		rc = 1
		for _, a := range s.as {
			if a != 0 {
				rc = 0
				break
			}
		}
	case funCheckBIsZero:
		rc = 1
		for _, b := range s.bs {
			if b != 0 {
				rc = 0
				break
			}
		}
	case funCheckAIsB:
		rc = 1
		for i := range s.as {
			if s.as[i] != s.bs[i] {
				rc = 0
				break
			}
		}
	case funSwapAB:
		for i := range s.as {
			s.as[i], s.bs[i] = s.bs[i], s.as[i]
		}
	case funAOrWithB:
		for i := range s.as {
			s.as[i] |= s.bs[i]
		}
	case funBOrWithA:
		for i := range s.as {
			s.bs[i] |= s.as[i]
		}
	case funAAndWithB:
		for i := range s.as {
			s.as[i] &= s.bs[i]
		}
	case funBAndWithA:
		for i := range s.as {
			s.bs[i] &= s.as[i]
		}
	case funAXorWithB:
		for i := range s.as {
			s.as[i] ^= s.bs[i]
		}
	case funBXorWithA:
		for i := range s.as {
			s.bs[i] ^= s.as[i]
		}
	case funAddAToB:
		var i big.Int
		bigA, bigB := s.getBigAB()
		i.Add(bigA, bigB)
		s.setBs(bigIntToPaddedBuffer(&i))
	case funAddBToA:
		var i big.Int
		bigA, bigB := s.getBigAB()
		i.Add(bigA, bigB)
		s.setAs(bigIntToPaddedBuffer(&i))
	case funSubAFromB:
		var i big.Int
		bigA, bigB := s.getBigAB()
		i.Sub(bigB, bigA)
		s.setBs(bigIntToPaddedBuffer(&i))
	case funSubBFromA:
		var i big.Int
		bigA, bigB := s.getBigAB()
		i.Sub(bigA, bigB)
		s.setAs(bigIntToPaddedBuffer(&i))
	case funMulAByB:
		var i big.Int
		bigA, bigB := s.getBigAB()
		i.Mul(bigA, bigB)
		s.setAs(bigIntToPaddedBuffer(&i))
	case funMulBByA:
		var i big.Int
		bigA, bigB := s.getBigAB()
		i.Mul(bigA, bigB)
		s.setBs(bigIntToPaddedBuffer(&i))
	case funDivAByB:
		var i big.Int
		bigA, bigB := s.getBigAB()
		if bigB.Cmp(math.BigZero) == 0 {
			return errCodeInvalidCode
		}
		i.Div(bigA, bigB)
		s.setAs(bigIntToPaddedBuffer(&i))
	case funDivBByA:
		var i big.Int
		bigA, bigB := s.getBigAB()
		if bigA.Cmp(math.BigZero) == 0 {
			return errCodeInvalidCode
		}
		i.Div(bigB, bigA)
		s.setBs(bigIntToPaddedBuffer(&i))
	default:
		if _, exists := funData[funNum]; exists {
			rc = getFunctionData(funNum)
		}
	}

	return rc
}

func (s *stateMachine) fun1(funNum int32, value int64) int64 {
	var rc int64

	switch { // echo
	case funNum == 1:
		fmt.Printf("%d\n", value)
	case funNum == 2:
		rc = value * 2 // double it
	case funNum == 3:
		rc = value / 2 // halve it
	case funNum == 26 || funNum == 33: // pay balance (prio:
		fmt.Printf("payout %d to account: %d", balance, value)
		balance = 0
	case funNum == 0x0110:
		s.setA(0, value)
	case funNum == 0x0111:
		s.setA(1, value)
	case funNum == 0x0112:
		s.setA(2, value)
	case funNum == 0x0113:
		s.setA(3, value)
	case funNum == 0x0116:
		s.setB(0, value)
	case funNum == 0x0117:
		s.setB(1, value)
	case funNum == 0x0118:
		s.setB(2, value)
	case funNum == 0x0119:
		s.setB(3, value)
	default:
		if _, exists := funData[funNum]; exists {
			rc = getFunctionData(funNum)
		}
	}

	return rc
}

func (s *stateMachine) fun2(funNum int32, value1, value2 int64) int64 {
	var rc int64

	switch {
	case funNum == 2:
		rc = value1 * value2 // multiply values
	case funNum == 3:
		rc = value1 / value2 // divide values
	case funNum == 4:
		rc = value1 + value2 // sum values
	case funNum == 31: // send amount to address
		if value1 > balance {
			value1 = balance
		}
		fmt.Printf("payout %d to account: %d", value1, value2)
		balance -= value1
	case funNum == 0x0114: // Set_A1_A2
		s.setA(0, value1)
		s.setA(1, value2)
	case funNum == 0x0115: // Set_A3_A4
		s.setA(2, value1)
		s.setA(3, value2)
	case funNum == 0x011a: // Set_B1_B2
		s.setB(0, value1)
		s.setB(1, value2)
	case funNum == 0x011b: // Set_B3_B4
		s.setB(2, value1)
		s.setB(3, value2)
	default:
		if _, exists := funData[funNum]; exists {
			rc = getFunctionData(funNum)
		}
	}

	return rc
}

func (s *stateMachine) invalidAddr(addr *int32) bool {
	return *addr < 0 || *addr > maxToMultiply || *addr*8 < 0 || *addr*8+8 > s.dsize
}

// generics :'(
func (s *stateMachine) invalidAddr64(addr *int64) bool {
	return *addr < 0 || *addr > maxToMultiply || *addr*8 < 0 || *addr*8+8 > int64(s.dsize)
}

func (s *stateMachine) getFun(fun *int16) int32 {
	if s.pc+2 >= s.csize {
		return -errCodeOverflow
	}

	*fun = *(*int16)(unsafe.Pointer(&s.code[s.pc+1]))

	return ok
}

func (s *stateMachine) getAddr(addr *int32) int32 {
	if s.pc+4 >= s.csize {
		return errCodeOverflow
	}

	*addr = *(*int32)(unsafe.Pointer(&s.code[s.pc+1]))

	if s.invalidAddr(addr) {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) getAddrWithOffset(off int32, addr *int32) int32 {
	if s.pc+4 >= s.csize {
		return errCodeOverflow
	}

	*addr = *(*int32)(unsafe.Pointer(&s.code[s.pc+1+off]))

	if s.invalidAddr(addr) {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) getAddrs(addr1, addr2 *int32) int32 {
	if s.pc+4 >= s.csize {
		return errCodeOverflow
	}

	*addr1 = *(*int32)(unsafe.Pointer(&s.code[s.pc+1]))
	*addr2 = *(*int32)(unsafe.Pointer(&s.code[s.pc+1+4]))

	if s.invalidAddr(addr1) || s.invalidAddr(addr2) {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) getAddrOff(addr *int32, off *int8) int32 {
	if s.pc+4+1 >= s.csize {
		return errCodeOverflow
	}

	*addr = *(*int32)(unsafe.Pointer(&s.code[s.pc+1]))
	*off = int8(s.code[s.pc+1+4])

	if s.invalidAddr(addr) || s.pc+int32(*off) >= s.csize {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) getAddrsOff(addr1, addr2 *int32, off *int8) int32 {
	if s.pc+4+4+1 >= s.csize {
		return errCodeOverflow
	}

	*addr1 = *(*int32)(unsafe.Pointer(&s.code[s.pc+1]))
	*addr2 = *(*int32)(unsafe.Pointer(&s.code[s.pc+1+4]))
	*off = int8(s.code[s.pc+1+4+4])

	if s.invalidAddr(addr1) || s.invalidAddr(addr2) || s.pc+int32(*off) >= s.csize {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) getFunAddr(fun *int16, addr *int32) int32 {
	if s.pc+2+4 >= s.csize {
		return errCodeOverflow
	}

	*fun = *(*int16)(unsafe.Pointer(&s.code[s.pc+1]))
	*addr = *(*int32)(unsafe.Pointer(&s.code[s.pc+1+2]))

	if s.invalidAddr(addr) {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) getFunAddrs(fun *int16, addr1, addr2 *int32) int32 {
	if s.pc+2+4+4 >= s.csize {
		return errCodeOverflow
	}

	*fun = *(*int16)(unsafe.Pointer(&s.code[s.pc+1]))
	*addr1 = *(*int32)(unsafe.Pointer(&s.code[s.pc+1+2]))
	*addr2 = *(*int32)(unsafe.Pointer(&s.code[s.pc+1+2+4]))

	if s.invalidAddr(addr1) || s.invalidAddr(addr2) {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) getAddrVal(addr *int32, val *int64) int32 {
	if s.pc+4+8 >= s.csize {
		return errCodeOverflow
	}

	*addr = *(*int32)(unsafe.Pointer(&s.code[s.pc+1]))
	*val = *(*int64)(unsafe.Pointer(&s.code[s.pc+1+4]))

	if s.invalidAddr(addr) {
		return errCodeOverflow
	}

	return ok
}

func (s *stateMachine) processOp(cssize, ussize int32) int32 {
	var rc int32

	if s.csize < 1 || s.pc >= s.csize {
		return ok
	}

	if s.pc < 0 {
		return errCodeInvalidCode
	}

	op := s.code[s.pc]

	switch op {
	case opCodeNop:
		for {
			rc++
			s.pc++
			if s.pc >= cssize || s.code[s.pc] != opCodeNop {
				break
			}
		}
	case opCodeSetVal:
		var addr int32
		var val int64
		rc = s.getAddrVal(&addr, &val)
		if rc == ok {
			rc = 1 + 4 + 8
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr*8])) = val
		}
	case opCodeSetDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) += *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
		}
	case opCodeClrDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr*8])) = 0
		}
	case opCodeIncDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr*8]))++
		}
	case opCodeDecDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr*8]))--
		}
	case opCodeNotDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr*8])) = ^*(*int64)(unsafe.Pointer(&s.data[addr*8]))
		}
	case opCodeAddDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) += *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
		}
	case opCodeSubDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) -= *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
		}
	case opCodeMulDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) *= *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
		}
	case opCodeDivDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			val := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
			if val == 0 {
				return errCodeInvalidCode
			}
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) /= val
		}
	case opCodeBorDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) |= *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
		}
	case opCodeAndDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) &= *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
		}
	case opCodeXorDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) ^= *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
		}
	case opCodeSetInd:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			// code uses a int64, but this should work without 64bit...
			addr := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
			if s.invalidAddr64(&addr) {
				rc = errCodeOverflow
			} else {
				s.pc += rc
				*(*int64)(unsafe.Pointer(&s.data[addr1*8])) =
					*(*int64)(unsafe.Pointer(&s.data[addr*8]))
			}
		}
	case opCodeSetIdx:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		var size int32 = 4 + 4
		if rc == ok {
			var addr3 int32
			rc = s.getAddrWithOffset(size, &addr3)
			if rc == ok {
				rc = 1 + size + 4

				base := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
				offs := *(*int64)(unsafe.Pointer(&s.data[addr3*8]))

				addr := base + offs
				if s.invalidAddr64(&addr) {
					rc = errCodeOverflow
				} else {
					s.pc += rc
					*(*int64)(unsafe.Pointer(&s.data[addr1*8])) =
						*(*int64)(unsafe.Pointer(&s.data[addr*8]))
				}
			}
		}
	case opCodePshDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			if s.us == ussize/8 {
				rc = errCodeOverflow
			} else {
				s.pc += rc
				s.us++
				*(*int64)(unsafe.Pointer(&s.data[s.dsize+cssize+ussize-s.us*8])) =
					*(*int64)(unsafe.Pointer(&s.data[addr*8]))
			}
		}
	case opCodePopDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			if s.us == 0 {
				rc = errCodeOverflow
			} else {
				s.pc += rc
				*(*int64)(unsafe.Pointer(&s.data[addr*8])) =
					*(*int64)(unsafe.Pointer(&s.data[s.dsize+cssize+ussize-s.us*8]))
				s.us--
			}
		}
	case opCodeJmpSub:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			if s.cs == cssize/8 {
				rc = errCodeOverflow
			}
			s.cs++
			*(*int64)(unsafe.Pointer(&s.data[s.dsize+cssize-s.cs*8])) = int64(s.pc) + int64(rc)
			s.pc = addr
		}
	case opCodeRetSub:
		if s.cs == 0 {
			rc = errCodeOverflow
		} else {
			rc = 1

			val := *(*int64)(unsafe.Pointer(&s.data[s.dsize+cssize-s.cs*8]))
			addr := int32(val)
			s.cs--
			s.pc = addr
		}
	case opCodeIndDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4

			addr := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))

			if s.invalidAddr64(&addr) {
				rc = errCodeOverflow
			} else {
				s.pc += rc
				*(*int64)(unsafe.Pointer(&s.data[addr*8])) =
					*(*int64)(unsafe.Pointer(&s.data[addr2*8]))
			}
		}
	case opCodeIdxDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		var size int32 = 4 + 4
		if rc == ok {
			var addr3 int32
			rc = s.getAddrWithOffset(size, &addr3)
			if rc == ok {
				rc = 1 + size + 4

				base := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
				offs := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

				addr := base + offs
				if s.invalidAddr64(&addr) {
					rc = errCodeOverflow
				} else {
					s.pc += rc
					*(*int64)(unsafe.Pointer(&s.data[addr*8])) =
						*(*int64)(unsafe.Pointer(&s.data[addr3*8]))
				}
			}
		}
	case opCodeModDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr2*8])) %= *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
		}
	case opCodeShlDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr2*8])) <<= *(*uint64)(unsafe.Pointer(&s.data[addr1*8]))
		}
	case opCodeShrDat:
		var addr1, addr2 int32
		rc = s.getAddrs(&addr1, &addr2)
		if rc == ok {
			rc = 1 + 4 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr2*8])) >>= *(*uint64)(unsafe.Pointer(&s.data[addr1*8]))
		}
	case opCodeJmpAdr:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			s.pc = addr
		}
	case opCodeBzrDat:
		var off int8
		var addr int32
		rc = s.getAddrOff(&addr, &off)
		if rc == ok {
			rc = 1 + 4 + 1
			val := *(*int64)(unsafe.Pointer(&s.data[addr*8]))
			if val == 0 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeBnzDat:
		var off int8
		var addr int32
		rc = s.getAddrOff(&addr, &off)
		if rc == ok {
			rc = 1 + 4 + 1
			val := *(*int64)(unsafe.Pointer(&s.data[addr*8]))
			if val != 0 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeBgtDat:
		var off int8
		var addr1, addr2 int32
		rc = s.getAddrsOff(&addr1, &addr2, &off)
		if rc == ok {
			rc = 1 + 4 + 4 + 1

			val1 := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

			if val1 > val2 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeBltDat:
		var off int8
		var addr1, addr2 int32
		rc = s.getAddrsOff(&addr1, &addr2, &off)
		if rc == ok {
			rc = 1 + 4 + 4 + 1

			val1 := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

			if val1 < val2 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeBgeDat:
		var off int8
		var addr1, addr2 int32
		rc = s.getAddrsOff(&addr1, &addr2, &off)
		if rc == ok {
			rc = 1 + 4 + 4 + 1

			val1 := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

			if val1 >= val2 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeBleDat:
		var off int8
		var addr1, addr2 int32
		rc = s.getAddrsOff(&addr1, &addr2, &off)
		if rc == ok {
			rc = 1 + 4 + 4 + 1

			val1 := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

			if val1 <= val2 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeBeqDat:
		var off int8
		var addr1, addr2 int32
		rc = s.getAddrsOff(&addr1, &addr2, &off)
		if rc == ok {
			rc = 1 + 4 + 4 + 1

			val1 := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

			if val1 == val2 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeBneDat:
		var off int8
		var addr1, addr2 int32
		rc = s.getAddrsOff(&addr1, &addr2, &off)
		if rc == ok {
			rc = 1 + 4 + 4 + 1

			val1 := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

			if val1 != val2 {
				s.pc += int32(off)
			} else {
				s.pc += rc
			}
		}
	case opCodeSlpDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			// NOTE: The "sleep_until" state value would be set to the current block + $addr.
			s.pc += rc
		}
	case opCodeFizDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			val := *(*int64)(unsafe.Pointer(&s.data[addr*8]))
			if val == 0 {
				s.pc = s.pcs
				s.finished = true
			} else {
				rc = 1 + 4
				s.pc += rc
			}
		}
	case opCodeStzDat:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			val := *(*int64)(unsafe.Pointer(&s.data[addr*8]))
			if val == 0 {
				s.pc += rc
				s.stopped = true
			} else {
				rc = 1 + 4
				s.pc += rc
			}
		}
	case opCodeFinImd:
		s.pc = s.pcs
		s.finished = true
	case opCodeStpImd:
		s.pc = s.pcs
		s.stopped = true
	case opCodeSlpImd:
		if rc == ok {
			rc = 1
			s.pc += rc
		}
	case opCodeErrAdr:
		var addr int32
		rc = s.getAddr(&addr)
		if rc == ok {
			rc = 1 + 4
			s.pce = addr
			rc = errCodeUnexpectedError
		}
	case opCodeSetPcs:
		rc = 1
		s.pc += rc
		s.pcs = s.pc
	case opCodeExtFun:
		var funNum int16
		rc = s.getFun(&funNum)
		if rc == ok {
			rc = 1 + 2
			s.pc += rc
			s.fun(int32(funNum))
		}
	case opCodeExtFunDat:
		var funNum int16
		var addr int32
		rc = s.getFunAddr(&funNum, &addr)
		if rc == ok {
			rc = 1 + 2 + 4
			s.pc += rc
			val := *(*int64)(unsafe.Pointer(&s.data[addr*8]))
			s.fun1(int32(funNum), val)
		}
	case opCodeExtFunDat2:
		var funNum int16
		var addr1, addr2 int32
		rc = s.getFunAddrs(&funNum, &addr1, &addr2)
		if rc == ok {
			rc = 1 + 2 + 4 + 4
			s.pc += rc
			val1 := *(*int64)(unsafe.Pointer(&s.data[addr1*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
			s.fun2(int32(funNum), val1, val2)
		}
	case opCodeExtFunRet:
		var funNum int16
		var addr int32
		rc = s.getFunAddr(&funNum, &addr)
		if rc == ok {
			rc = 1 + 2 + 4
			s.pc += rc
			*(*int64)(unsafe.Pointer(&s.data[addr*8])) = s.fun(int32(funNum))
		}
	case opCodeExtFunRetDat:
		var funNum int16
		var addr1, addr2 int32
		rc = s.getFunAddrs(&funNum, &addr1, &addr2)
		if rc == ok {
			rc = 1 + 2 + 4 + 4
			s.pc += rc

			val := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))

			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) = s.fun1(int32(funNum), val)
		}
	case opCodeExtFunRetDat2:
		var funNum int16
		var addr1, addr2, addr3 int32
		rc = s.getFunAddrs(&funNum, &addr1, &addr2)
		if rc == ok {
			rc = s.getAddrWithOffset(2+4+4, &addr3)
		}
		if rc == ok {
			rc = 1 + 2 + 4 + 4 + 4
			s.pc += rc

			val := *(*int64)(unsafe.Pointer(&s.data[addr2*8]))
			val2 := *(*int64)(unsafe.Pointer(&s.data[addr3*8]))

			*(*int64)(unsafe.Pointer(&s.data[addr1*8])) = s.fun2(int32(funNum), val, val2)
		}
	default:
		rc = errCodeInvalidCode
	}

	if rc == errCodeOverflow && s.pce != 0 {
		rc = ok
		s.pc = s.pce
	}

	if rc >= ok {
		s.steps++
	}

	return rc
}

func (s *stateMachine) execute() error {
	for {
		rc := s.processOp(int32(callStackPages*callStackePageBytes), int32(userStackPages*userStackPageBytes))

		if rc >= ok {
			if s.stopped {
				return nil
			} else if s.finished {
				return nil
			}
		} else {
			switch rc {
			case errCodeOverflow:
				return ErrOverflow
			case errCodeInvalidCode:
				return ErrInvalidCode
			default:
				return ErrUnexpectedError
			}
		}
	}
}
