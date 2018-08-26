package at

import (
	"errors"
	"fmt"
	"unsafe"
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

	a1 int64
	a2 int64
	a3 int64
	a4 int64

	b1 int64
	b2 int64
	b3 int64
	b4 int64

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

func (s *stateMachine) fun(funNum int32) int64 {
	var rc int64

	switch {
	case funNum == 1:
		rc = s.val
	case funNum == 2:
		if s.val == 9 {
			rc = 0
			s.val = 0
		} else {
			s.val++
			rc = s.val
		}
	case funNum == 3: // get size
		rc = 10
	case funNum == 4:
		rc = int64(funNum)
	case funNum == 25 || funNum == 32:
		rc = balance
		if _, exists := funData[funNum]; exists {
			for _, f := range funData {
				f.offset = 0
			}
		}
	case funNum == 0x0100:
		rc = s.a1
	case funNum == 0x0101:
		rc = s.a2
	case funNum == 0x0102:
		rc = s.a3
	case funNum == 0x0103:
		rc = s.a4
	case funNum == 0x0104:
		rc = s.b1
	case funNum == 0x0105:
		rc = s.b2
	case funNum == 0x0106:
		rc = s.b3
	case funNum == 0x0107:
		rc = s.b4
	case funNum == 0x0120:
		s.a1 = 0
		s.a2 = 0
		s.a3 = 0
		s.a4 = 0
	case funNum == 0x0121:
		s.b1 = 0
		s.b2 = 0
		s.b3 = 0
		s.b4 = 0
	case funNum == 0x0122:
		s.a1 = 0
		s.a2 = 0
		s.a3 = 0
		s.a4 = 0

		s.b1 = 0
		s.b2 = 0
		s.b3 = 0
		s.b4 = 0
	case funNum == 0x0123:
		s.a1 = s.b1
		s.a2 = s.b2
		s.a3 = s.b3
		s.a4 = s.b4
	case funNum == 0x0124:
		s.b1 = s.a1
		s.b2 = s.a2
		s.b3 = s.a3
		s.b4 = s.a4
	case funNum == 0x0125: // Check_A_Is_Zero
		if s.a1 == 0 && s.a2 == 0 && s.a3 == 0 && s.a4 == 0 {
			rc = 1
		}
	case funNum == 0x0126: // Check_B_Is_Zero
		if s.b1 == 0 && s.b2 == 0 && s.b3 == 0 && s.b4 == 0 {
			rc = 1
		}
	case funNum == 0x0127: // Check_A_Equals_B
		if s.a1 == s.b1 && s.a2 == s.b2 && s.a3 == s.b3 && s.a4 == s.b4 {
			rc = 1
		}
	case funNum == 0x0128: // Swap_A_and_B
		s.a1, s.b1 = s.b1, s.a1
		s.a2, s.b2 = s.b2, s.a2
		s.a3, s.b3 = s.b3, s.a3
		s.a4, s.b4 = s.b4, s.a4
	case funNum == 0x0129: // OR_A_with_B
		s.a1 = s.a1 | s.b1
		s.a2 = s.a2 | s.b2
		s.a3 = s.a3 | s.b3
		s.a4 = s.a4 | s.b4
	case funNum == 0x012a: // OR_B_with_A
		s.b1 = s.a1 | s.b1
		s.b2 = s.a2 | s.b2
		s.b3 = s.a3 | s.b3
		s.b4 = s.a4 | s.b4
	case funNum == 0x012b: // AND_A_with_B
		s.a1 = s.a1 & s.b1
		s.a2 = s.a2 & s.b2
		s.a3 = s.a3 & s.b3
		s.a4 = s.a4 & s.b4
	case funNum == 0x012c: // AND_B_with_A
		s.b1 = s.a1 & s.b1
		s.b2 = s.a2 & s.b2
		s.b3 = s.a3 & s.b3
		s.b4 = s.a4 & s.b4
	case funNum == 0x012d: // XOR_A_with_B
		s.a1 = s.a1 ^ s.b1
		s.a2 = s.a2 ^ s.b2
		s.a3 = s.a3 ^ s.b3
		s.a4 = s.a4 ^ s.b4
	case funNum == 0x012e: // XOR_B_with_A
		s.b1 = s.a1 ^ s.b1
		s.b2 = s.a2 ^ s.b2
		s.b3 = s.a3 ^ s.b3
		s.b4 = s.a4 ^ s.b4
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
		s.a1 = value
	case funNum == 0x0111:
		s.a2 = value
	case funNum == 0x0112:
		s.a3 = value
	case funNum == 0x0113:
		s.a4 = value
	case funNum == 0x0116:
		s.b1 = value
	case funNum == 0x0117:
		s.b2 = value
	case funNum == 0x0118:
		s.b3 = value
	case funNum == 0x0119:
		s.b4 = value
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
		s.a1 = value1
		s.a2 = value2
	case funNum == 0x0115: // Set_A3_A4
		s.a3 = value1
		s.a4 = value2
	case funNum == 0x011a: // Set_B1_B2
		s.b1 = value1
		s.b2 = value2
	case funNum == 0x011b: // Set_B3_B4
		s.b3 = value1
		s.b4 = value2
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
				return -2
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
