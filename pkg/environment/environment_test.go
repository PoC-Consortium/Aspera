package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimumFee(t *testing.T) {
	assert.Equal(t, uint64(0), MinimumFee(0))
	assert.Equal(t, uint64(oneBurst), MinimumFee(preDymaxionForkHeight-1))
	assert.Equal(t, uint64(feeQuant), MinimumFee(preDymaxionForkHeight))
}

func TestBlockReward(t *testing.T) {
	assert.Equal(t, int64(0*oneBurst), BlockReward(0))
	assert.Equal(t, int64(0*oneBurst), BlockReward(1944000))
	assert.Equal(t, int64(10000*oneBurst), BlockReward(10799))
	assert.Equal(t, int64(9500*oneBurst), BlockReward(10800))
	assert.Equal(t, int64(1498*oneBurst), BlockReward(400000))
}
