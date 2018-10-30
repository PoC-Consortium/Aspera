package environment

const (
	oneBurst = 100000000
	feeQuant = 735000

	preDymaxionForkHeight = 500000
	poc2ForkHeight        = 502000
)

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
