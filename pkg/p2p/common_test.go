package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlignMilestonesWithCurrent(t *testing.T) {
	milestones := []*blockMeta{
		&blockMeta{height: 0},
	}
	rs := alignMilestonesWithCurrent(
		milestones,
		&blockMeta{height: 0},
	)
	assert.EqualValues(t, []*blockRange{&blockRange{from: *milestones[0]}}, rs)

	milestones = []*blockMeta{
		&blockMeta{height: 0},
		&blockMeta{height: 10},
		&blockMeta{height: 20},
	}
	rs = alignMilestonesWithCurrent(
		milestones,
		&blockMeta{height: 0},
	)
	assert.EqualValues(
		t,
		[]*blockRange{
			&blockRange{from: *milestones[0], to: milestones[1]},
			&blockRange{from: *milestones[1], to: milestones[2]},
			&blockRange{from: *milestones[2]},
		},
		rs,
	)

	rs = alignMilestonesWithCurrent(
		milestones,
		&blockMeta{height: 5},
	)
	assert.EqualValues(
		t,
		[]*blockRange{
			&blockRange{from: blockMeta{height: 5}, to: milestones[1]},
			&blockRange{from: *milestones[1], to: milestones[2]},
			&blockRange{from: *milestones[2]},
		},
		rs,
	)

	milestones = []*blockMeta{
		&blockMeta{height: 0},
		&blockMeta{height: 10},
		&blockMeta{height: 20},
	}
	rs = alignMilestonesWithCurrent(
		milestones,
		&blockMeta{height: 10},
	)
	assert.EqualValues(
		t,
		[]*blockRange{
			&blockRange{from: *milestones[1], to: milestones[2]},
			&blockRange{from: *milestones[2]},
		},
		rs,
	)
}
