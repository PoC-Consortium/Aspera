package shabal256

import (
	"fmt"
	"testing"
)

type shabal256Test struct {
	out string
	in  string
}

var shabal256Tests = []shabal256Test{
	{"AEC750D11FEEE9F16271922FBAF5A9BE142F62019EF8D720F858940070889014", ""},
	{"6E99A12CEFC40F9ED1C5B049BAB4FF04F2F81315CA42AF49AE2450BBCED92870", "a"},
	{"08CCF5B9E70B0EC06F29FA719729331FEFA1831AA05B10DE6418D2C446782B79", "ab"},
	{"07225FAB83CA48FB480D22219410D5CA008359EFBFD315829029AFE2CB3F0404", "abc"},
	{"1589E078F1B22B2500ECFE905A4DF6624C48D65E43F2FFC7BDCF129A2B8D95E5", "abcd"},
	{"B1E1E43CAD2C32DC2D0EEBDBC7FFDBBD935B780082CF48BB3A3FD9867FEDF410", "abcde"},
	{"E905102C7A89BB19FAF68EFE5511167D07EFF63E0E5BD8E1CC1C866E3983E9C3", "abcdef"},
	{"4E90A530B3258B2A02964C46AA7BEC8A7D867806E6B9C2F5A66D3F55C529F1CF", "abcdefg"},
	{"27F1E7C6B1D9049E6E6C39DEC8086A7E0BF371E35276C3319BBA3C1177A5C6C3", "abcdefgh"},
	{"4730607C763946D3EAFF5B17A347BCC94C1512B5FB3823BD08E2B3B3AE79CBC4", "abcdefghi"},
	{"54D67E2444797C7A5FDD14A86431C3DB02D9874F6EF66C9A2EBD9D07076F3670", "abcdefghij"},
	{"FAF2B841B4E39BFA753259DF67E307C9B98010A5FF848D379C67D45241D97203", "Discard medicine more than two years old."},
	{"AAAAF35B7EB0FFEB00EBFE7AADBEAF83228FD755FDEB752F1741B87CF9DD95C2", "He who has a shady past knows that nice guys finish last."},
	{"FA39027CD6290A50D1F919D453E7E892107BE3749A03061A0ADA5726855D11C9", "I wouldn't marry him with a ten foot pole."},
	{"8C32B7523BFCCB775FB68F2B5F3860067161C40906E86C4A32EA275A31747148", "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{"AFD302415CD1B8D139D9D63195CD08CAA096B85A03CE10F2AE8318FAB38EBC34", "The days of the digital watch are numbered.  -Tom Stoppard"},
	{"B412255E525F5E0D3F38F4F305C35BE3E2E2AB88C8EF14C2ADB8F3D698BC7BA7", "Nepal premier won't resign."},
	{"D116E846D73FF429D5B75EC5E66D0B3C526A0BB7FC0BFD89C7B6BCE3FB03DAD4", "For every action there is an equal and opposite government program."},
	{"2FE392498C645FC3F6E72635A468C6158107B055C8DD527266B1DAE942BF1CDD", "His money is twice tainted: 'taint yours and 'taint mine."},
	{"889F8C9975B20FFC1D3A77D7B7E5CFE69ADB5F47728337936BC8435D53903C5D", "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{"34407A29F0846CD18FF8362E97050B1A5DABE3185E937629A56064A179383280", "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{"466D403C817663C785754ED4F84EC7CB0159186D8EDF0B79A21A3155D35FFC2C", "size:  a.out:  bad magic"},
	{"3A7113CFF9BC27485DDAAED8C3107FD5E877568AEB5E32E1076328E7619BFB81", "The major problem is with sendmail.  -Mark Horton"},
	{"41869346A0D42CB6899E4D3AFF5E5E957457429861344684B8706594A55F32B7", "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{"596D90772030E6D74012FBDF648BEC736942D8BCA6FBD838C3F028DDD920A892", "If the enemy is within range, then so are you."},
	{"40C3C7C952480869C8A447653F74815F05CA72E71287AF0804B497AC2BBDEC53", "It's well we cannot hear the screams/That we create in others' dreams."},
	{"8B113E74B5DAB12ABCA060802889247511EB25702846B9EAB8E76F5824D27C53", "You remind me of a TV show, but that's all right: I watch it anyway."},
	{"7A199D44DE2BA8338AFA374A510FD62A2F93DD32FB5869CB6BF44E07FC40CAE7", "C is as portable as Stonehedge!!"},
	{"6D32F5FA6053C4B83D954D96DB2CD6D36697F7A1A305BF4F7C59DAD1FC868935", "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{"6505CEF17C67468B447DBC6DA21E0219E83B7EBF882D5FA09A3A239DA0932F56", "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{"894B6F60B07224C3A43A9FFB2842ECD9CF564DF782CE4F1BCC5BE24476B6C404", "How can you write a big system without C++?  -Paul Glick"},
}

func TestSize(t *testing.T) {
	c := New()
	if got := c.Size(); got != Size {
		t.Errorf("Size = %d; want %d", got, Size)
	}
}

func TestBlockSize(t *testing.T) {
	c := New()
	if got := c.BlockSize(); got != BlockSize {
		t.Errorf("BlockSize = %d want %d", got, BlockSize)
	}
}

func TestSum2556(t *testing.T) {
	for _, test := range shabal256Tests {
		s := fmt.Sprintf("%X", Sum256([]byte(test.in)))
		if s != test.out {
			t.Fatalf("Sum256 function: shabal256(%s) = %s want %s", test.in, s, test.out)
		}
	}
}

var bench = New()
var buf = make([]byte, 8192)

func benchmarkSize(b *testing.B, size int) {
	b.SetBytes(int64(size))
	sum := make([]byte, bench.Size())
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:size])
		bench.Sum(sum[:0])
	}
}

func BenchmarkHash8Bytes(b *testing.B) {
	benchmarkSize(b, 8)
}

func BenchmarkHash1K(b *testing.B) {
	benchmarkSize(b, 1024)
}

func BenchmarkHash8K(b *testing.B) {
	benchmarkSize(b, 8192)
}
