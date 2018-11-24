package block

import (
	"testing"

	api "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/stretchr/testify/assert"
)

func TestSetBaseTargetAndCummulativeDifficulty(t *testing.T) {
	// TODO: test cummulative difficutly
	type test struct {
		block                    *Block
		previousBlocks           []*Block
		expCommulativeDifficulty []byte
		expBaseTarget            uint64
	}

	tests := []test{
		test{
			block:          &Block{Block: &api.Block{}},
			previousBlocks: nil,
			expBaseTarget:  18325193796,
		},
		test{
			block: &Block{Block: &api.Block{Height: 1}},
			previousBlocks: []*Block{
				&Block{Block: &api.Block{
					Height: 0,
				}},
			},
			expBaseTarget: 18325193796,
		},
		test{
			block: &Block{Block: &api.Block{Height: 3}},
			previousBlocks: []*Block{
				&Block{Block: &api.Block{
					Height: 0,
				}},
				&Block{Block: &api.Block{
					Height: 1,
				}},
				&Block{Block: &api.Block{
					Height: 2,
				}},
			},
			expBaseTarget: 18325193796,
		},
		test{
			block: &Block{Block: &api.Block{
				Height:    4,
				Timestamp: 1088,
			}},
			previousBlocks: []*Block{
				&Block{Block: &api.Block{
					Height:     0,
					BaseTarget: 18325193796,
					Timestamp:  0,
				}},
				&Block{Block: &api.Block{
					Height:     1,
					BaseTarget: 18325193796,
					Timestamp:  683,
				}},
				&Block{Block: &api.Block{
					Height:     2,
					BaseTarget: 18325193796,
					Timestamp:  795,
				}},
				&Block{Block: &api.Block{
					Height:     3,
					BaseTarget: 18325193796,
					Timestamp:  966,
				}},
			},
			expBaseTarget: 18325193796,
		},
		test{
			block: &Block{Block: &api.Block{
				Height:    2699,
				Timestamp: 697716,
			}},
			previousBlocks: []*Block{
				&Block{Block: &api.Block{
					Height:     2675,
					BaseTarget: 90361595,
					Timestamp:  691977,
				}},
				&Block{Block: &api.Block{
					Height:     2676,
					BaseTarget: 88963419,
					Timestamp:  692205,
				}},
				&Block{Block: &api.Block{
					Height:     2677,
					BaseTarget: 82688043,
					Timestamp:  692283,
				}},
				&Block{Block: &api.Block{
					Height:     2678,
					BaseTarget: 87610512,
					Timestamp:  692573,
				}},
				&Block{Block: &api.Block{
					Height:     2679,
					BaseTarget: 96146481,
					Timestamp:  693171,
				}},
				&Block{Block: &api.Block{
					Height:     2680,
					BaseTarget: 97737324,
					Timestamp:  694046,
				}},
				&Block{Block: &api.Block{
					Height:     2681,
					BaseTarget: 100150149,
					Timestamp:  694203,
				}},
				&Block{Block: &api.Block{
					Height:     2682,
					BaseTarget: 104952227,
					Timestamp:  694378,
				}},
				&Block{Block: &api.Block{
					Height:     2683,
					BaseTarget: 109721199,
					Timestamp:  694444,
				}},
				&Block{Block: &api.Block{
					Height:     2684,
					BaseTarget: 92826201,
					Timestamp:  694474,
				}},
				&Block{Block: &api.Block{
					Height:     2685,
					BaseTarget: 91721199,
					Timestamp:  694577,
				}},
				&Block{Block: &api.Block{
					Height:     2686,
					BaseTarget: 89824685,
					Timestamp:  694634,
				}},
				&Block{Block: &api.Block{
					Height:     2687,
					BaseTarget: 86420988,
					Timestamp:  694776,
				}},
				&Block{Block: &api.Block{
					Height:     2688,
					BaseTarget: 81178441,
					Timestamp:  694899,
				}},
				&Block{Block: &api.Block{
					Height:     2689,
					BaseTarget: 78557695,
					Timestamp:  695290,
				}},
				&Block{Block: &api.Block{
					Height:     2690,
					BaseTarget: 75770897,
					Timestamp:  695500,
				}},
				&Block{Block: &api.Block{
					Height:     2691,
					BaseTarget: 72433804,
					Timestamp:  695555,
				}},
				&Block{Block: &api.Block{
					Height:     2692,
					BaseTarget: 69286688,
					Timestamp:  695603,
				}},
				&Block{Block: &api.Block{
					Height:     2693,
					BaseTarget: 81413498,
					Timestamp:  696631,
				}},
				&Block{Block: &api.Block{
					Height:     2694,
					BaseTarget: 82198843,
					Timestamp:  696657,
				}},
				&Block{Block: &api.Block{
					Height:     2695,
					BaseTarget: 83966528,
					Timestamp:  696719,
				}},
				&Block{Block: &api.Block{
					Height:     2696,
					BaseTarget: 87138027,
					Timestamp:  697081,
				}},
				&Block{Block: &api.Block{
					Height:     2697,
					BaseTarget: 75311301,
					Timestamp:  697189,
				}},
				&Block{Block: &api.Block{
					Height:     2698,
					BaseTarget: 73938306,
					Timestamp:  697449,
				}},
			},
			expBaseTarget: 83175285,
		},
		test{
			block: &Block{Block: &api.Block{
				Height:    2700,
				Timestamp: 697772,
			}},
			previousBlocks: []*Block{
				&Block{Block: &api.Block{
					Height:     2676,
					BaseTarget: 88963419,
					Timestamp:  692205,
				}},
				&Block{Block: &api.Block{
					Height:     2677,
					BaseTarget: 82688043,
					Timestamp:  692283,
				}},
				&Block{Block: &api.Block{
					Height:     2678,
					BaseTarget: 87610512,
					Timestamp:  692573,
				}},
				&Block{Block: &api.Block{
					Height:     2679,
					BaseTarget: 96146481,
					Timestamp:  693171,
				}},
				&Block{Block: &api.Block{
					Height:     2680,
					BaseTarget: 97737324,
					Timestamp:  694046,
				}},
				&Block{Block: &api.Block{
					Height:     2681,
					BaseTarget: 100150149,
					Timestamp:  694203,
				}},
				&Block{Block: &api.Block{
					Height:     2682,
					BaseTarget: 104952227,
					Timestamp:  694378,
				}},
				&Block{Block: &api.Block{
					Height:     2683,
					BaseTarget: 109721199,
					Timestamp:  694444,
				}},
				&Block{Block: &api.Block{
					Height:     2684,
					BaseTarget: 92826201,
					Timestamp:  694474,
				}},
				&Block{Block: &api.Block{
					Height:     2685,
					BaseTarget: 91721199,
					Timestamp:  694577,
				}},
				&Block{Block: &api.Block{
					Height:     2686,
					BaseTarget: 89824685,
					Timestamp:  694634,
				}},
				&Block{Block: &api.Block{
					Height:     2687,
					BaseTarget: 86420988,
					Timestamp:  694776,
				}},
				&Block{Block: &api.Block{
					Height:     2688,
					BaseTarget: 81178441,
					Timestamp:  694899,
				}},
				&Block{Block: &api.Block{
					Height:     2689,
					BaseTarget: 78557695,
					Timestamp:  695290,
				}},
				&Block{Block: &api.Block{
					Height:     2690,
					BaseTarget: 75770897,
					Timestamp:  695500,
				}},
				&Block{Block: &api.Block{
					Height:     2691,
					BaseTarget: 72433804,
					Timestamp:  695555,
				}},
				&Block{Block: &api.Block{
					Height:     2692,
					BaseTarget: 69286688,
					Timestamp:  695603,
				}},
				&Block{Block: &api.Block{
					Height:     2693,
					BaseTarget: 81413498,
					Timestamp:  696631,
				}},
				&Block{Block: &api.Block{
					Height:     2694,
					BaseTarget: 82198843,
					Timestamp:  696657,
				}},
				&Block{Block: &api.Block{
					Height:     2695,
					BaseTarget: 83966528,
					Timestamp:  696719,
				}},
				&Block{Block: &api.Block{
					Height:     2696,
					BaseTarget: 87138027,
					Timestamp:  697081,
				}},
				&Block{Block: &api.Block{
					Height:     2697,
					BaseTarget: 75311301,
					Timestamp:  697189,
				}},
				&Block{Block: &api.Block{
					Height:     2698,
					BaseTarget: 73938306,
					Timestamp:  697449,
				}},
				&Block{Block: &api.Block{
					Height:     2699,
					BaseTarget: 83175285,
					Timestamp:  697716,
				}},
			},
			expBaseTarget: 83362225,
		},
	}
	for _, test := range tests {
		b := test.block
		b.SetBaseTargetAndCummulativeDifficulty(test.previousBlocks)
		assert.Equal(t, test.expBaseTarget, b.BaseTarget)
	}
}
