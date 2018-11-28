package blockchain

import (
	"fmt"
)

func uint64ToBs(i uint64) []byte {
	return []byte(fmt.Sprintf("%020d", i))
}
