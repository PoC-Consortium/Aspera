package environment

import (
	"math/big"
)

const (
	oneBurst = 100000000
	feeQuant = 735000

	preDymaxionForkHeight = 500000
	poc2ForkHeight        = 502000

	RewardRecipientStartHeight = 6500

	AdjustDifficutlyHeight = 2700

	InitialBaseTarget = 18325193796
	MaxBaseTarget     = 18325193796
)

// Minimum returns the minimum transaction fee at the specified height.
func MinimumFee(height int32) uint64 {
	switch {
	case height == 0:
		return 0
	case height < preDymaxionForkHeight:
		return oneBurst
	default:
		return feeQuant
	}
}

// BlockReward returns the block reward at the specified height.
func BlockReward(height int32) int64 {
	if height == 0 || height >= 1944000 {
		return 0
	}
	month := int64(height / 10800)
	reward := big.NewInt(10000)
	tmp1 := big.NewInt(95)
	tmp1.Exp(tmp1, big.NewInt(month), nil)
	reward.Mul(reward, tmp1)
	tmp2 := big.NewInt(100)
	tmp2.Exp(tmp2, big.NewInt(month), nil)
	reward.Quo(reward, tmp2)
	return reward.Int64() * oneBurst
}
