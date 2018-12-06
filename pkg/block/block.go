package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"

	api "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	. "github.com/PoC-Consortium/Aspera/pkg/blockchain"
	"github.com/PoC-Consortium/Aspera/pkg/burstmath"
	. "github.com/PoC-Consortium/Aspera/pkg/common/math"
	"github.com/PoC-Consortium/Aspera/pkg/crypto"
	"github.com/PoC-Consortium/Aspera/pkg/crypto/shabal256"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
	env "github.com/PoC-Consortium/Aspera/pkg/environment"
	"github.com/PoC-Consortium/Aspera/pkg/transaction"

	"github.com/dgraph-io/badger"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
)

var (
	ErrInvalidBlockVersion         = errors.New("invalid block version")
	ErrBlockUnexpectedLen          = errors.New("block unexpected length in byte serialisation")
	ErrPreviousBlockMismatch       = errors.New("previous block id doesn't match current block's")
	ErrTimestampTooLate            = errors.New("timestamp to late")
	ErrTimestampSmallerPrevious    = errors.New("timestamp smaller than previous block's")
	ErrGenerationSignatureMismatch = errors.New("generation signature mismatch")
	ErrInvalidPayloadHash          = errors.New("invalid payload hash")
	ErrBlockAmountTooLow           = errors.New("block's total amount too low for transactions")
	ErrBlockFeeTooLow              = errors.New("blocks' total fee too low for transactions")
)

const (
	generationSignatureLen = 64
	// TODO: move constants
	oneBurst               = 100000000
	maxTimestampDifference = 15

	ById     = "block:id:%020d"
	ByHeight = "block:height:%010d"
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

func FromProtoBytes(bs []byte) (*Block, error) {
	pbBlock := new(api.Block)
	if err := proto.Unmarshal(bs, pbBlock); err != nil {
		return nil, err
	}
	return NewBlock(pbBlock)
}

func (b *Block) CalcScoop() uint32 {
	return burstmath.CalcScoop(b.Height, b.GenerationSignature)
}

func (b *Block) ToBytes() []byte {
	e := encoding.NewEncoder(b.SizeInBytes())

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
	l := 4 + 4 + 8 + 4 + 4 + 32 + 32 + 32 + 8 + 64 + 32
	if b.Version < 3 {
		l += 4 + 4
	} else {
		l += 8 + 8
	}
	if b.BlockATs != nil {
		l += len(b.BlockATs)
	}
	return l
}

func (b *Block) CalculateHashAndID() ([32]byte, uint64) {
	return crypto.BytesToHashAndID(b.ToBytes())
}

func (b *Block) PreValidate(previous *Block) error {
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

	generationSignatureExp := CalculateGenerationSignature(previous)
	for i := range b.GenerationSignature {
		if generationSignatureExp[i] != b.GenerationSignature[i] {
			return ErrGenerationSignatureMismatch
		}
	}

	if len(b.transactions) > 0 {
		if err := b.validateTransactions(now); err != nil {
			return err
		}
	}

	_, b.Id = b.CalculateHashAndID()

	return nil
}

func (b *Block) Validate(previousBlocks []*Block) error {
	b.SetBaseTargetAndCumulativeDifficulty(previousBlocks)
	// TODO: sequential validation
	return nil
}

func (b *Block) BlockReward() int64 {
	blockReward := env.BlockReward(b.Height)
	for _, tx := range b.transactions {
		blockReward += int64(tx.GetHeader().Fee)
	}
	return blockReward
}

func (b *Block) Execute() error {
	if b.Height < env.RewardRecipientStartHeight {
		if err := BC.TransferBurst(nil, b.GeneratorID(), b.BlockReward(), 0); err != nil {
			panic(err)
		}
	} else {
		panic("process reward reipient")
	}
	for _, tx := range b.transactions {
		if err := transaction.Execute(tx); err != nil {
			panic(err)
		}
	}
	return nil
}

func (b *Block) validateTransactions(now uint32) error {
	payloadDigest := sha256.New()
	var totalFee, totalAmount int64
	for _, tx := range b.transactions {
		if bs, err := transaction.ValidateAndGetBytes(tx, b.Height, b.Timestamp, now); err == nil {
			payloadDigest.Write(bs)
		} else {
			return err
		}
		h := tx.GetHeader()
		// TODO: either we use int64 for all amounts or uint64
		totalFee += int64(h.Fee)
		totalAmount += int64(h.Amount)
	}

	switch {
	case totalAmount > b.TotalAmount:
		return ErrBlockAmountTooLow
	case totalFee > b.TotalFee:
		return ErrBlockFeeTooLow
	case !bytes.Equal(payloadDigest.Sum(nil), b.PayloadHash):
		return ErrInvalidPayloadHash
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

func (b *Block) GeneratorID() uint64 {
	_, id := crypto.BytesToHashAndID(b.GeneratorPublicKey)
	return id
}

func (b *Block) SetBaseTargetAndCumulativeDifficulty(previousBlocks []*Block) {
	switch {
	case b.Height == 0:
		b.BaseTarget = env.InitialBaseTarget
		b.CumulativeDifficulty = big.NewInt(0).Bytes()
	case b.Height < 4:
		b.BaseTarget = env.InitialBaseTarget
		previousBlock := previousBlocks[len(previousBlocks)-1]
		cumulativeDifficulty := new(big.Int).SetBytes(previousBlock.CumulativeDifficulty)
		var tmp big.Int
		tmp.Quo(MaxBig64, big.NewInt(env.InitialBaseTarget))
		b.CumulativeDifficulty = cumulativeDifficulty.Add(cumulativeDifficulty, &tmp).Bytes()
	case b.Height < env.AdjustDifficutlyHeight:
		var avgBaseTargetBig big.Int
		previousBlocks = previousBlocks[len(previousBlocks)-4:]
		for _, p := range previousBlocks {
			avgBaseTargetBig.Add(&avgBaseTargetBig, BigFromUint64(p.BaseTarget))
		}
		avgBaseTargetBig.Quo(&avgBaseTargetBig, big.NewInt(4))

		dt := int64(b.Timestamp - previousBlocks[0].Timestamp)

		currentBaseTarget := avgBaseTargetBig.Uint64()
		newBaseTargetBig := avgBaseTargetBig.Mul(&avgBaseTargetBig, big.NewInt(dt))
		newBaseTarget := newBaseTargetBig.Quo(newBaseTargetBig, big.NewInt(240*4)).Uint64()
		if newBaseTarget < 0 || newBaseTarget > env.MaxBaseTarget {
			newBaseTarget = env.MaxBaseTarget
		}
		if newBaseTarget < currentBaseTarget*9/10 {
			newBaseTarget = currentBaseTarget * 9 / 10
		}
		if newBaseTarget == 0 {
			newBaseTarget = 1
		}
		twofoldCurrentBaseTarget := int64(currentBaseTarget) * 11 / 10
		if twofoldCurrentBaseTarget < 0 {
			twofoldCurrentBaseTarget = env.MaxBaseTarget
		}
		if newBaseTarget > uint64(twofoldCurrentBaseTarget) {
			newBaseTarget = uint64(twofoldCurrentBaseTarget)
		}
		b.BaseTarget = newBaseTarget
		previousCumulativeDifficulty := new(big.Int).SetBytes(previousBlocks[3].CumulativeDifficulty)
		var tmp big.Int
		tmp.Quo(MaxBig64, BigFromUint64(newBaseTarget))
		b.CumulativeDifficulty = previousCumulativeDifficulty.Add(
			previousCumulativeDifficulty, &tmp).Bytes()
	default:
		previousBlocks = previousBlocks[len(previousBlocks)-24:]
		avgBaseTargetBig := BigFromUint64(previousBlocks[23].BaseTarget)
		for i := 22; i >= 0; i-- {
			avgBaseTargetBig.Mul(avgBaseTargetBig, big.NewInt(int64(24-i)))
			avgBaseTargetBig.Add(avgBaseTargetBig, BigFromUint64(previousBlocks[i].BaseTarget))
			avgBaseTargetBig.Quo(avgBaseTargetBig, big.NewInt(int64(25-i)))
		}
		dt := int64(b.Timestamp - previousBlocks[0].Timestamp)
		var dtTarget int64 = 24 * 4 * 60

		if dt < dtTarget/2 {
			dt = dtTarget / 2
		}
		if dt > dtTarget*2 {
			dt = dtTarget * 2
		}
		currentBaseTarget := previousBlocks[23].BaseTarget
		tmp1 := new(big.Int).Mul(avgBaseTargetBig, big.NewInt(dt))
		newBaseTarget := tmp1.Quo(tmp1, big.NewInt(dtTarget)).Uint64()
		if newBaseTarget < 0 || newBaseTarget > env.MaxBaseTarget {
			newBaseTarget = env.MaxBaseTarget
		}
		if newBaseTarget == 0 {
			newBaseTarget = 1
		}
		if newBaseTarget < currentBaseTarget*8/10 {
			newBaseTarget = currentBaseTarget * 8 / 10
		}
		if newBaseTarget > currentBaseTarget*12/10 {
			newBaseTarget = currentBaseTarget * 12 / 10
		}
		b.BaseTarget = newBaseTarget
		previousCumulativeDifficulty := new(big.Int).SetBytes(previousBlocks[23].CumulativeDifficulty)
		var tmp2 big.Int
		tmp2.Quo(MaxBig64, BigFromUint64(newBaseTarget))
		b.CumulativeDifficulty = previousCumulativeDifficulty.Add(
			previousCumulativeDifficulty, &tmp2).Bytes()
	}
}

func (b *Block) Freeze() ([][][]byte, error) {
	block := b.Block

	if pb, err := proto.Marshal(block); err != nil {
		return nil, err
	} else {
		return [][][]byte{
			[][]byte{
				[]byte(fmt.Sprintf(ById, block.Id)),
				[]byte(fmt.Sprintf(ByHeight, block.Height)),
			},
			[][]byte{
				[]byte(fmt.Sprintf(ByHeight, block.Height)),
				pb,
			},
		}, nil
	}
}

func Thaw(txn *badger.Txn, key string) (*Block, error) {
	keyBs := []byte(key)
	if ok, _ := regexp.Match("^block:(?:id|height):[0-9]+$", keyBs); !ok {
		return nil, errors.New("invalid key " + key)
	}
fetchData:
	if blockItem, err := txn.Get(keyBs); err != nil {
		return nil, err
	} else {
		if blockBs, err := blockItem.Value(); err != nil {
			return nil, err
		} else {
			if !bytes.HasPrefix(keyBs, []byte("block:height:")) {
				keyBs = blockBs
				goto fetchData
			}

			pbBlock := new(api.Block)
			if err := proto.Unmarshal(blockBs, pbBlock); err != nil {
				return nil, err
			}
			if block, err := NewBlock(pbBlock); err == nil {
				return block, nil
			} else {
				return nil, err
			}
		}
	}
}

func BulkThaw(txn *badger.Txn, key string, limit int) ([]*Block, error) {
	if block, err := Thaw(txn, key); err != nil {
		return nil, err
	} else {
		var blocks []*Block
		for i := 1; i <= limit; i++ {
			block, err = Thaw(txn, fmt.Sprintf(ByHeight, block.Height+1))
			if err == nil {
				blocks = append(blocks, block)
			} else if err == badger.ErrKeyNotFound {
				return blocks, nil
			} else {
				return nil, err
			}
		}
		return blocks, nil
	}
}
