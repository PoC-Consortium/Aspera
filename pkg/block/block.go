package block

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	api "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/burstmath"
	"github.com/ac0v/aspera/pkg/crypto"
	"github.com/ac0v/aspera/pkg/crypto/shabal256"
	"github.com/ac0v/aspera/pkg/encoding"
	"github.com/ac0v/aspera/pkg/transaction"

	"github.com/json-iterator/go"
)

var (
	ErrInvalidBlockVersion         = errors.New("invalid block version")
	ErrBlockUnexpectedLen          = errors.New("block unexpected length in byte serialisation")
	ErrPreviousBlockMismatch       = errors.New("previous block id doesn't match current block's")
	ErrTimestampTooLate            = errors.New("timestamp to late")
	ErrTimestampSmallerPrevious    = errors.New("timestamp smaller than previous block's")
	ErrGenerationSignatureMismatch = errors.New("generation signature mismatch")
)

const (
	generationSignatureLen = 64
	// TODO: move constants
	oneBurst               = 100000000
	maxTimestampDifference = 15
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Block struct {
	*api.Block
	transactions []transaction.Transaction
}

func NewBlock(apiBlock *api.Block) (*Block, error) {
	transactions := make([]transaction.Transaction, len(apiBlock.Transactions))
	for i, a := range apiBlock.Transactions {
		if tx, err := transaction.AnyToTransaction(a); err == nil {
			transactions[i] = tx
		} else {
			return nil, err
		}
	}
	return &Block{apiBlock, transactions}, nil
}

func (b *Block) CalcScoop() uint32 {
	return burstmath.CalcScoop(b.Height, b.GenerationSignature)
}

func (b *Block) ToBytes() []byte {
	e := encoding.NewEncoder(make([]byte, b.SizeInBytes()))

	e.WriteInt32(b.Version)
	e.WriteUint32(b.Timestamp)
	e.WriteUint64(b.PreviousBlock)
	e.WriteUint32(uint32(len(b.Transactions)))
	if b.Version < 3 {
		e.WriteInt32(int32(b.TotalAmount / oneBurst))
		e.WriteInt32(int32(b.TotalFee / oneBurst))
	} else {
		e.WriteInt64(b.TotalAmount)
		e.WriteInt64(b.TotalFee)
	}
	e.WriteUint32(b.PayloadLength)
	e.WriteBytes(b.PayloadHash)
	e.WriteBytes(b.GeneratorPublicKey)
	e.WriteBytes(b.GenerationSignature)
	if b.Version > 1 {
		e.WriteBytes(b.PreviousBlockHash)
	}
	e.WriteUint64(b.Nonce)
	if b.BlockATs != nil {
		e.WriteBytes(b.BlockATs)
	}
	e.WriteBytes(b.BlockSignature)

	return e.Bytes()
}

func (b *Block) SizeInBytes() int {
	l := 4 + 4 + 8 + 4 + 4 + 32 + 32 + 32 + 8 + 64
	if b.Version < 3 {
		l += 4 + 4
	} else {
		l += 8 + 8
	}
	if b.Version > 1 {
		l += 32
	}
	if b.BlockATs != nil {
		l += len(b.BlockATs)
	}
	return l
}

func (b *Block) CalculateHashAndID() ([32]byte, uint64) {
	return crypto.BytesToHashAndID(b.ToBytes())
}

func (b *Block) Validate(previous *Block) error {
	now := burstmath.DateToTimestamp(time.Now())

	switch {
	case b.Version != 3:
		return ErrInvalidBlockVersion
	case b.Timestamp <= previous.Timestamp:
		return ErrTimestampSmallerPrevious
	case b.Timestamp > now+maxTimestampDifference:
		return ErrTimestampTooLate
	}

	previousHash, previousID := previous.CalculateHashAndID()
	switch {
	case previousID != b.PreviousBlock:
		return ErrPreviousBlockMismatch
	case !bytes.Equal(previousHash[:], b.PreviousBlockHash):
		return ErrPreviousBlockMismatch
	}

	// ToDo: check for duplicte blocks - may this should go to the raw storage stuff
	// throw new BlockNotAcceptedException("Duplicate block or invalid id for block " + block.getHeight());

	// for _, t := range b.Transactions {
	// 	if err := t.VerifySignature(); err != nil {
	// 		return err
	// 	}
	// }

	generationSignatureExp := CalculateGenerationSignature(previous)
	for i := range b.GenerationSignature {
		if generationSignatureExp[i] != b.GenerationSignature[i] {
			return ErrGenerationSignatureMismatch
		}
	}

	return nil
}

func CalculateGenerationSignature(previous *Block) []byte {
	bs := make([]byte, 8)
	_, id := crypto.BytesToHashAndID(previous.GeneratorPublicKey)
	binary.BigEndian.PutUint64(bs, id)
	hash := shabal256.Sum256(append(previous.GenerationSignature, bs...))
	return hash[:]
}
