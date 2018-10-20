package attachment

var escrowDeadlineActionNameOf = map[uint8]string{
	0: "undecided",
	1: "release",
	2: "refund",
	3: "split",
}
var escrowDeadlineActionIdOf = map[string]uint8{}

func init() {
	for id, name := range escrowDeadlineActionNameOf {
		escrowDeadlineActionIdOf[name] = id
	}
}
