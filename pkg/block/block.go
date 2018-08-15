package block

import (
	"encoding/hex"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	"github.com/ac0v/aspera/pkg/burstmath"
)

type Block struct {
	*pb.Block
}

func NewBlock(b *pb.Block) *Block {
	return &Block{b}
}

func (b *Block) CalcScoop() uint32 {
	genSig, err := hex.DecodeString(b.GenerationSignature)
	if err != nil {
		panic("generation signature wrong format")
	}
	return burstmath.CalcScoop(b.Height, genSig)
}
